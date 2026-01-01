package app

import "time"

type Config struct {
	// SyncInterval is the interval between automatic syncs from providers.
	// If zero, background syncing is disabled.
	SyncInterval time.Duration
}
