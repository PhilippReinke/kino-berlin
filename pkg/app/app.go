package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
)

type App struct {
	storage   domain.Storage
	providers []domain.Provider

	// sync management
	syncInterval time.Duration
	syncCtx      context.Context
	syncCancel   context.CancelFunc
	syncWg       sync.WaitGroup
	syncMu       sync.RWMutex
	syncRunning  bool
}

func New(storage domain.Storage, providers []domain.Provider, config Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		storage:      storage,
		providers:    providers,
		syncInterval: config.SyncInterval,
		syncCtx:      ctx,
		syncCancel:   cancel,
	}
}

func (a *App) FetchScreenings(filters ...domain.Filter) ([]domain.Screening, error) {
	if a.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	screenings, err := a.storage.Fetch(filters...)
	if err != nil {
		return nil, fmt.Errorf("fetching screenings: %w", err)
	}

	return screenings, nil
}

func (a *App) GetAvailableCinemas() ([]string, error) {
	screenings, err := a.FetchScreenings(domain.ExpiredScreeningFilter())
	if err != nil {
		return nil, err
	}

	cinemaMap := make(map[string]bool)
	for _, s := range screenings {
		if s.Cinema != "" {
			cinemaMap[s.Cinema] = true
		}
	}

	cinemas := make([]string, 0, len(cinemaMap))
	for cinema := range cinemaMap {
		cinemas = append(cinemas, cinema)
	}

	return cinemas, nil
}

func (a *App) GetAvailableDates() ([]time.Time, error) {
	screenings, err := a.FetchScreenings(domain.ExpiredScreeningFilter())
	if err != nil {
		return nil, err
	}

	dateMap := make(map[string]time.Time)
	for _, s := range screenings {
		year, month, day := s.Start.UTC().Date()
		dateKey := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dateMap[dateKey] = date
	}

	dates := make([]time.Time, 0, len(dateMap))
	for _, date := range dateMap {
		dates = append(dates, date)
	}

	// sort dates (oldest first)
	for i := 0; i < len(dates)-1; i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i].After(dates[j]) {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	return dates, nil
}
