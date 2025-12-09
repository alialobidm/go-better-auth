package domain

import (
	"time"
)

// KeyValueStore represents the persistent key-value store table in the database.
// This is a domain model used for secondary storage operations.
type KeyValueStore struct {
	Key       string     `json:"key" gorm:"primaryKey"`
	Value     string     `json:"value"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// SecondaryStorageMemoryConfig holds settings specific to the in-memory storage.
type SecondaryStorageMemoryConfig struct {
	// CleanupInterval controls how often expired entries are cleaned up.
	// If zero, the implementation should use a default.
	CleanupInterval time.Duration
}

// SecondaryStorageDatabaseConfig holds settings specific to the database storage.
type SecondaryStorageDatabaseConfig struct {
	// CleanupInterval controls how often expired entries are cleaned up.
	// If zero, the implementation should use a default.
	CleanupInterval time.Duration
}
