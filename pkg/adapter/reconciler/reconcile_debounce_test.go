package reconciler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingletonDebounce_AllowOrMark(t *testing.T) {
	t.Run("first call should allow", func(t *testing.T) {
		// given
		debounce := &SingletonDebounce{}
		window := 100 * time.Millisecond
		// when
		allowed := debounce.AllowOrMark(window)
		// then
		assert.True(t, allowed)
	})

	t.Run("immediate second call should mark pending", func(t *testing.T) {
		// given
		debounce := &SingletonDebounce{}
		window := 100 * time.Millisecond
		allowed := debounce.AllowOrMark(window)
		assert.True(t, allowed)
		// when
		allowed = debounce.AllowOrMark(window)
		// then
		assert.False(t, allowed)
		assert.True(t, debounce.pending)
	})

	t.Run("call after window should allow", func(t *testing.T) {
		// given
		debounce := &SingletonDebounce{}
		window := 100 * time.Millisecond
		allowed := debounce.AllowOrMark(window)
		assert.True(t, allowed)
		// when
		time.Sleep(window + 10*time.Millisecond) // wait for window to pass
		allowed = debounce.AllowOrMark(window)
		// then
		assert.True(t, allowed)
		assert.False(t, debounce.pending)
	})
}

func TestSingletonDebounce_ShouldRequeue(t *testing.T) {
	t.Run("no pending should not requeue", func(t *testing.T) {
		// given
		debounce := &SingletonDebounce{}
		// when
		shouldRequeue, remaining := debounce.ShouldRequeue()
		// then
		assert.False(t, shouldRequeue)
		assert.Equal(t, time.Duration(0), remaining)
	})

	t.Run("with pending should requeue", func(t *testing.T) {
		// given
		debounce := &SingletonDebounce{}
		window := 100 * time.Millisecond
		// Mark as pending
		debounce.AllowOrMark(window)
		debounce.AllowOrMark(window) // This should mark pending
		// when
		shouldRequeue, remaining := debounce.ShouldRequeue()
		// then
		assert.True(t, shouldRequeue)
		assert.True(t, remaining > 0)
		assert.False(t, debounce.pending) // Should clear pending
	})

	t.Run("after window expired should requeue with no remaining time", func(t *testing.T) {
		// given
		debounce := &SingletonDebounce{}
		window := 100 * time.Millisecond
		// Mark as pending and wait
		debounce.AllowOrMark(window)
		debounce.AllowOrMark(window) // This should mark pending
		time.Sleep(window + 10*time.Millisecond)
		// when
		shouldRequeue, remaining := debounce.ShouldRequeue()
		// then
		assert.True(t, shouldRequeue)
		assert.Equal(t, time.Duration(0), remaining)
		assert.False(t, debounce.pending) // Should clear pending
	})
}

func TestSingletonDebounce_Concurrency(t *testing.T) {
	debounce := &SingletonDebounce{}
	window := 50 * time.Millisecond

	// Test that multiple goroutines can safely access the debounce
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			debounce.AllowOrMark(window)
			debounce.ShouldRequeue()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("Test timed out - possible deadlock")
		}
	}
}
