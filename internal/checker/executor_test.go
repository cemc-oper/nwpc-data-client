package checker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandExecutorProcessesInOrder(t *testing.T) {
	// onTaskDone runs in the consumer goroutine after each result is
	// processed, so recording there observes the execution order.
	var hooked []time.Duration

	exec := newCommandExecutor(3, nil, time.Now(), "", func(result CheckResult) {
		hooked = append(hooked, result.ForecastTime)
	})

	results := []CheckResult{
		{ForecastTime: 0, FilePath: "/a"},
		{ForecastTime: 3 * time.Hour, FilePath: "/b"},
		{ForecastTime: 6 * time.Hour, FilePath: "/c"},
	}
	for _, r := range results {
		exec.submit(r)
	}
	exec.finish()
	exec.wait()

	require.Len(t, hooked, 3)
	assert.Equal(t, []time.Duration{0, 3 * time.Hour, 6 * time.Hour}, hooked)
}

func TestCommandExecutorSubmitDoesNotBlock(t *testing.T) {
	// Capacity covers every submission, so producers never block even
	// before the consumer has drained anything.
	exec := newCommandExecutor(2, nil, time.Now(), "", nil)
	done := make(chan struct{})
	go func() {
		defer close(done)
		exec.submit(CheckResult{ForecastTime: 0})
		exec.submit(CheckResult{ForecastTime: time.Hour})
		exec.finish()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("submit blocked despite sufficient queue capacity")
	}
	exec.wait()
}
