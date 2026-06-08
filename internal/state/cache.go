package state

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	State      string
	LastActive time.Time
}

type StateCache struct {
	users map[string]entry
	mu    sync.RWMutex
}

var activeCache = StateCache{
	users: make(map[string]entry),
}

func GetState(phone string) (string, bool) {
	activeCache.mu.RLock()
	defer activeCache.mu.RUnlock()
	record, exists := activeCache.users[phone]
	return record.State, exists
}

func SetState(phone, state string) {
	activeCache.mu.Lock()
	defer activeCache.mu.Unlock()
	activeCache.users[phone] = entry{
		State:      state,
		LastActive: time.Now(),
	}
}

func StartSweeper(timeout time.Duration, ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				activeCache.mu.Lock()
				for phone, record := range activeCache.users {
					if time.Since(record.LastActive) > timeout {
						delete(activeCache.users, phone)
					}
				}
				activeCache.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}

	}()
}
