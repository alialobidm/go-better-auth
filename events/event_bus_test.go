package events

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/GoBetterAuth/go-better-auth/config"
	"github.com/GoBetterAuth/go-better-auth/models"
)

func getMockConfig() *models.Config {
	return config.NewConfig()
}

func TestEventBus_Publish(t *testing.T) {
	bus := NewEventBus(getMockConfig(), nil)
	defer bus.Close()

	payload, _ := json.Marshal(map[string]string{
		"user_id": "123",
	})

	event := models.Event{
		Type:      models.EventUserSignedUp,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
		Metadata: map[string]string{
			"source": "test",
		},
	}
	err := bus.Publish(context.Background(), event)
	assert.NoError(t, err)
}

func TestWatermillEventBus_Publish(t *testing.T) {
	bus := NewEventBus(getMockConfig(), nil)
	defer bus.Close()

	payload, _ := json.Marshal(map[string]string{
		"user_id": "123",
	})

	event := models.Event{
		Type:      models.EventUserSignedUp,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
		Metadata: map[string]string{
			"source": "test",
		},
	}
	err := bus.Publish(context.Background(), event)
	assert.NoError(t, err)
}

func TestWatermillEventBus_Subscribe(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := NewEventBus(getMockConfig(), nil)
		defer bus.Close()

		var wg sync.WaitGroup
		handlerCalled := atomic.Bool{}
		var receivedEvent models.Event

		wg.Add(1)
		_, err := bus.Subscribe(models.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
			handlerCalled.Store(true)
			receivedEvent = event
			wg.Done()
			return nil
		})
		assert.NoError(t, err)

		payload, _ := json.Marshal(map[string]string{
			"user_id": "456",
		})

		event := models.Event{
			Type:      models.EventUserSignedUp,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
			Metadata:  map[string]string{"source": "test"},
		}
		err = bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		wg.Wait()
		assert.True(t, handlerCalled.Load())
		assert.Equal(t, models.EventUserSignedUp, receivedEvent.Type)
	})
}

func TestWatermillEventBus_MultipleEvents(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := NewEventBus(getMockConfig(), nil)
		defer bus.Close()

		signupCount := atomic.Int32{}
		loginCount := atomic.Int32{}
		var wg sync.WaitGroup

		wg.Add(2)

		_, err := bus.Subscribe(models.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
			signupCount.Add(1)
			wg.Done()
			return nil
		})
		assert.NoError(t, err)

		_, err = bus.Subscribe(models.EventUserLoggedIn, func(ctx context.Context, event models.Event) error {
			loginCount.Add(1)
			wg.Done()
			return nil
		})
		assert.NoError(t, err)

		signupPayload, _ := json.Marshal(map[string]string{"user_id": "789"})
		loginPayload, _ := json.Marshal(map[string]string{"user_id": "789"})

		signupEvent := models.Event{
			Type:      models.EventUserSignedUp,
			Timestamp: time.Now().UTC(),
			Payload:   signupPayload,
			Metadata:  map[string]string{"source": "test"},
		}
		loginEvent := models.Event{
			Type:      models.EventUserLoggedIn,
			Timestamp: time.Now().UTC(),
			Payload:   loginPayload,
			Metadata:  map[string]string{"source": "test"},
		}
		bus.Publish(context.Background(), signupEvent)
		bus.Publish(context.Background(), loginEvent)

		wg.Wait()
		assert.Greater(t, signupCount.Load(), int32(0))
		assert.Greater(t, loginCount.Load(), int32(0))
	})
}

func TestWatermillEventBus_EventData(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := NewEventBus(getMockConfig(), nil)
		defer bus.Close()

		var wg sync.WaitGroup
		var receivedPayload json.RawMessage

		wg.Add(1)
		_, err := bus.Subscribe(models.EventUserLoggedIn, func(ctx context.Context, event models.Event) error {
			receivedPayload = event.Payload
			wg.Done()
			return nil
		})
		assert.NoError(t, err)

		payload, _ := json.Marshal(map[string]string{
			"user_id":  "123",
			"username": "testuser",
			"email":    "test@example.com",
		})

		event := models.Event{
			Type:      models.EventUserLoggedIn,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
			Metadata:  map[string]string{"source": "test"},
		}
		err = bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		wg.Wait()
		assert.NotNil(t, receivedPayload)

		// Unmarshal to verify the data
		var data map[string]string
		json.Unmarshal(receivedPayload, &data)
		assert.Equal(t, "123", data["user_id"])
	})
}

func TestEventBus_WithCustomPubSub(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Test that EventBus works with a custom PubSub implementation
		config := getMockConfig()
		customPubSub := NewInMemoryPubSub()
		bus := NewEventBus(config, customPubSub)
		defer bus.Close()

		var wg sync.WaitGroup
		var receivedEvent models.Event

		wg.Add(1)
		_, err := bus.Subscribe(models.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
			receivedEvent = event
			wg.Done()
			return nil
		})
		assert.NoError(t, err)

		// Give subscription time to set up
		time.Sleep(10 * time.Millisecond)

		payload, _ := json.Marshal(map[string]string{
			"user_id": "custom-transport-test",
			"email":   "test@example.com",
		})

		event := models.Event{
			Type:      models.EventUserSignedUp,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
			Metadata: map[string]string{
				"source": "custom_pubsub_test",
			},
		}

		err = bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		wg.Wait()
		assert.Equal(t, models.EventUserSignedUp, receivedEvent.Type)

		// Unmarshal to verify the payload
		var data map[string]string
		json.Unmarshal(receivedEvent.Payload, &data)
		assert.Equal(t, "custom-transport-test", data["user_id"])
	})
}

func TestEventBus_MultipleHandlersPerTopic(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// This test verifies that multiple handlers subscribed to the same event type
		// all receive the event without creating multiple consumer goroutines
		bus := NewEventBus(getMockConfig(), nil)
		defer bus.Close()

		const numHandlers = 10
		counts := make([]atomic.Int32, numHandlers)
		var wg sync.WaitGroup
		wg.Add(numHandlers)

		// Subscribe multiple handlers to the same event type
		for i := range numHandlers {
			handlerIndex := i
			_, err := bus.Subscribe(models.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
				counts[handlerIndex].Add(1)
				wg.Done()
				return nil
			})
			assert.NoError(t, err)
		}

		payload, _ := json.Marshal(map[string]string{"user_id": "multi-handler-test"})
		event := models.Event{
			Type:      models.EventUserSignedUp,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
			Metadata:  map[string]string{"source": "test"},
		}

		// Publish one event
		err := bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		// All handlers should receive the same event
		wg.Wait()
		for i := range numHandlers {
			assert.Equal(t, int32(1), counts[i].Load(), "handler %d should have received exactly 1 event", i)
		}
	})
}

func TestEventBus_Unsubscribe(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := NewEventBus(getMockConfig(), nil)
		defer bus.Close()

		var wg sync.WaitGroup
		count := atomic.Int32{}

		wg.Add(1)
		// Subscribe first handler
		id1, err := bus.Subscribe(models.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
			count.Add(1)
			return nil
		})
		assert.NoError(t, err)

		// Subscribe second handler
		id2, err := bus.Subscribe(models.EventUserSignedUp, func(ctx context.Context, event models.Event) error {
			count.Add(1)
			wg.Done()
			return nil
		})
		assert.NoError(t, err)

		// Unsubscribe the first handler
		bus.Unsubscribe(models.EventUserSignedUp, id1)

		payload, _ := json.Marshal(map[string]string{"user_id": "unsubscribe-test"})
		event := models.Event{
			Type:      models.EventUserSignedUp,
			Timestamp: time.Now().UTC(),
			Payload:   payload,
			Metadata:  map[string]string{"source": "test"},
		}

		// Publish event
		err = bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		wg.Wait()
		// Only the second handler should have received the event
		assert.Equal(t, int32(1), count.Load())

		// Clean up: unsubscribe second handler
		bus.Unsubscribe(models.EventUserSignedUp, id2)
	})
}
