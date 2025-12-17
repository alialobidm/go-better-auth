package config

import (
	"github.com/GoBetterAuth/go-better-auth/models"
)

// NewPlugin creates a new plugin with the provided options. Useful for simple plugins.
func NewPlugin(options ...models.PluginOption) models.Plugin {
	// Create a new BasePlugin instance
	plugin := &models.BasePlugin{}
	// Set default values
	plugin.SetMetadata(models.PluginMetadata{})
	plugin.SetConfig(models.PluginConfig{
		Enabled: true,
	})
	plugin.SetInit(func(ctx *models.PluginContext) error {
		return nil
	})
	plugin.SetMigrations(make([]any, 0))
	plugin.SetRoutes(make([]models.PluginRoute, 0))
	plugin.SetRateLimit(&models.PluginRateLimit{})
	plugin.SetDatabaseHooks(&models.PluginDatabaseHooks{})
	plugin.SetEventHooks(&models.PluginEventHooks{})
	plugin.SetClose(func() error {
		return nil
	})

	for _, option := range options {
		option(plugin)
	}

	return plugin
}

func WithPluginMetadata(metadata models.PluginMetadata) models.PluginOption {
	return func(p models.Plugin) {
		p.SetMetadata(metadata)
	}
}

func WithPluginConfig(config models.PluginConfig) models.PluginOption {
	return func(p models.Plugin) {
		p.SetConfig(config)
	}
}

func WithPluginInit(init func(ctx *models.PluginContext) error) models.PluginOption {
	return func(p models.Plugin) {
		p.SetInit(init)
	}
}

func WithPluginMigrations(migrations []any) models.PluginOption {
	return func(p models.Plugin) {
		p.SetMigrations(migrations)
	}
}

func WithPluginRoutes(routes []models.PluginRoute) models.PluginOption {
	return func(p models.Plugin) {
		p.SetRoutes(routes)
	}
}

func WithPluginRateLimit(rateLimit *models.PluginRateLimit) models.PluginOption {
	return func(p models.Plugin) {
		p.SetRateLimit(rateLimit)
	}
}

func WithPluginDatabaseHooks(databaseHooks *models.PluginDatabaseHooks) models.PluginOption {
	return func(p models.Plugin) {
		p.SetDatabaseHooks(databaseHooks)
	}
}

func WithPluginEventHooks(eventHooks *models.PluginEventHooks) models.PluginOption {
	return func(p models.Plugin) {
		p.SetEventHooks(eventHooks)
	}
}

func WithPluginClose(close func() error) models.PluginOption {
	return func(p models.Plugin) {
		p.SetClose(close)
	}
}
