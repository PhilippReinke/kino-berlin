package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
)

func (a *App) StartBackgroundSync() error {
	if a.syncInterval == 0 {
		// sync disabled
		return nil
	}

	a.syncMu.Lock()
	defer a.syncMu.Unlock()

	if a.syncRunning {
		return fmt.Errorf("background sync already running")
	}

	a.syncRunning = true
	a.syncWg.Add(1)

	a.syncWg.Go(func() {
		defer a.syncWg.Done()
		ticker := time.NewTicker(a.syncInterval)
		defer ticker.Stop()

		// initial sync immediately
		log.Printf("Starting background sync (interval: %v)", a.syncInterval)
		if err := a.SyncFromProviders(a.syncCtx); err != nil {
			log.Printf("Initial background sync failed: %v", err)
		}

		for {
			select {
			case <-a.syncCtx.Done():
				log.Printf("Background sync stopped")
				return
			case <-ticker.C:
				log.Printf("Running scheduled sync")
				if err := a.SyncFromProviders(a.syncCtx); err != nil {
					log.Printf("Background sync failed: %v", err)
				}
			}
		}
	})

	return nil
}

func (a *App) SyncFromProviders(ctx context.Context) error {
	if len(a.providers) == 0 {
		return fmt.Errorf("no providers configured")
	}

	for _, provider := range a.providers {
		if err := a.syncFromProvider(ctx, provider); err != nil {
			log.Printf("Failed to sync from provider %q: %v", provider.Name(), err)
		}
	}

	return nil
}

func (a *App) syncFromProvider(ctx context.Context, provider domain.Provider) error {
	log.Printf("Start scraping %q.", provider.Name())

	screenings, err := provider.Scrape()
	if err != nil {
		return fmt.Errorf("scraping failed: %w", err)
	}

	for _, screening := range screenings {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := a.storage.Upsert(screening); err != nil {
			log.Printf("Failed to upsert screening %q: %v", screening.ID, err)
		}
	}

	log.Printf("Finished scraping %q.", provider.Name())

	return nil
}

func (a *App) StopBackgroundSync() {
	a.syncMu.Lock()
	if !a.syncRunning {
		a.syncMu.Unlock()
		return
	}
	a.syncMu.Unlock()

	a.syncCancel()
	a.syncWg.Wait()

	a.syncMu.Lock()
	a.syncRunning = false

	// new context for potential restart
	a.syncCtx, a.syncCancel = context.WithCancel(context.Background())
	a.syncMu.Unlock()
}
