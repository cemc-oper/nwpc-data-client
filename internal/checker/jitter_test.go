package checker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJitteredIntervalZeroJitter(t *testing.T) {
	interval := 5 * time.Second
	for i := 0; i < 100; i++ {
		assert.Equal(t, interval, jitteredInterval(interval, 0))
	}
}

func TestJitteredIntervalWithinBounds(t *testing.T) {
	interval := 10 * time.Second
	jitter := 0.3
	lower := time.Duration(float64(interval) * (1 - jitter))
	upper := time.Duration(float64(interval) * (1 + jitter))
	for i := 0; i < 1000; i++ {
		got := jitteredInterval(interval, jitter)
		assert.GreaterOrEqual(t, got, lower)
		assert.LessOrEqual(t, got, upper)
	}
}

func TestStartupJitterWithinBounds(t *testing.T) {
	interval := 5 * time.Second
	for i := 0; i < 1000; i++ {
		got := startupJitter(interval)
		assert.GreaterOrEqual(t, got, time.Duration(0))
		assert.Less(t, got, interval)
	}
}

func TestStartupJitterNonPositiveInterval(t *testing.T) {
	assert.Equal(t, time.Duration(0), startupJitter(0))
	assert.Equal(t, time.Duration(0), startupJitter(-time.Second))
}

func TestSleepContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- sleepContext(ctx, time.Hour)
	}()
	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("context-aware sleep did not return after cancellation")
	}
}

func TestSleepContextAlreadyCancelledWithZeroDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.ErrorIs(t, sleepContext(ctx, 0), context.Canceled)
}
