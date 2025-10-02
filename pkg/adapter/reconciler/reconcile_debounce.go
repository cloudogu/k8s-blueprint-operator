package reconciler

import (
	"sync"
	"time"
)

// SingletonDebounce coalesces bursts for a single blueprint.
type SingletonDebounce struct {
	mu      sync.Mutex
	next    time.Time // when the next event is allowed
	pending bool      // true if something arrived during cooldown
}

// AllowOrMark returns true if we should enqueue now.
// If we're still in cooldown, it marks pending and returns false.
func (d *SingletonDebounce) AllowOrMark(window time.Duration) bool {
	now := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if now.Before(d.next) {
		d.pending = true
		return false
	}
	d.next = now.Add(window)
	d.pending = false
	return true
}

// ShouldRequeue tells if we should do a trailing reconciliation after cooldown.
// Returns (shouldRequeue, remainingCooldown).
func (d *SingletonDebounce) ShouldRequeue() (bool, time.Duration) {
	now := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.pending {
		return false, 0
	}
	remaining := time.Duration(0)
	if now.Before(d.next) {
		remaining = d.next.Sub(now)
	}
	d.pending = false
	return true, remaining
}
