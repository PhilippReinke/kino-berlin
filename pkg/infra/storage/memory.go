package storage

import (
	"fmt"
	"sort"
	"sync"

	"github.com/PhilippReinke/kino-berlin/pkg/domain"
)

type Memory struct {
	mu         sync.RWMutex
	screenings map[domain.ScreeningID]domain.Screening
}

var _ domain.Storage = &Memory{}

func NewMemory() *Memory {
	return &Memory{
		mu:         sync.RWMutex{},
		screenings: make(map[domain.ScreeningID]domain.Screening),
	}
}

func (m *Memory) Upsert(s domain.Screening) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, ok := m.screenings[s.ID]
	if ok && current.UpdatedAt.After(s.UpdatedAt) {
		return fmt.Errorf("existing screening is newer")
	}

	m.screenings[s.ID] = s

	return nil
}

func (m *Memory) Fetch(filter ...domain.Filter) ([]domain.Screening, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var screenings []domain.Screening
	for _, s := range m.screenings {
		keep := true
		for _, f := range filter {
			if !f(s) {
				keep = false
				break
			}
		}

		if keep {
			screenings = append(screenings, s)
		}
	}

	// Sort screenings for deterministic order:
	// 1. By start time (earliest first)
	// 2. By cinema (alphabetically)
	// 3. By title (alphabetically)
	sort.Slice(screenings, func(i, j int) bool {
		if !screenings[i].Start.Equal(screenings[j].Start) {
			return screenings[i].Start.Before(screenings[j].Start)
		}
		if screenings[i].Cinema != screenings[j].Cinema {
			return screenings[i].Cinema < screenings[j].Cinema
		}
		return screenings[i].Title < screenings[j].Title
	})

	return screenings, nil
}
