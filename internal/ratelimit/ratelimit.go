package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alerts for the same port within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
// Events for the same key will be suppressed until the cooldown elapses.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if the event identified by key should be forwarded.
// It returns false if the same key was allowed within the cooldown window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded timestamp for key, allowing the next event
// through immediately regardless of the cooldown.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush removes all recorded timestamps, effectively resetting the limiter.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
