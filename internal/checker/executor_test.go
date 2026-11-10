package checker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandExecutorProcessesInOrder(t *testing.T) {
	// onTaskDone runs in the consumer goroutine after each result is
	// processed, so recording there observes the execution order.
	var hooked []time.Duration

	exec := newCommandExecutor(context.Background(), 0, 3, nil, time.Now(), "", func(result CheckResult) {
		hooked = append(hooked, result.ForecastTime)
	}, nil)

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
	exec := newCommandExecutor(context.Background(), 0, 2, nil, time.Now(), "", nil, nil)
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

func TestCommandExecutorReturnsCheckFailureAndStops(t *testing.T) {
	exec := newCommandExecutor(context.Background(), 0, 3, nil, time.Now(), "", nil, nil)
	exec.submit(CheckResult{ForecastTime: 0, FilePath: "/a"})
	exec.submit(CheckResult{ForecastTime: time.Hour, Error: assert.AnError})
	exec.submit(CheckResult{ForecastTime: 2 * time.Hour, FilePath: "/c"})
	exec.finish()

	err := exec.wait()

	require.ErrorIs(t, err, assert.AnError)
	assert.False(t, exec.submit(CheckResult{ForecastTime: 3 * time.Hour}))
}

func TestCommandExecutorReturnsCommandFailureAndStops(t *testing.T) {
	templates, err := buildCommandTemplates("exit 7", nil)
	require.NoError(t, err)
	exec := newCommandExecutor(context.Background(), 0, 2, templates, time.Now(), "", nil, nil)
	exec.submit(CheckResult{ForecastTime: 0, FilePath: "/a"})
	exec.submit(CheckResult{ForecastTime: time.Hour, FilePath: "/b"})
	exec.finish()

	err = exec.wait()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exit code is not 0: 7")
}

func TestCommandExecutorReturnsCommandTimeoutAndStops(t *testing.T) {
	templates, err := buildCommandTemplates(commandHelperCommand(t, "wait"), nil)
	require.NoError(t, err)
	exec := newCommandExecutor(context.Background(), 50*time.Millisecond, 2, templates, time.Now(), "", nil, nil)
	exec.submit(CheckResult{ForecastTime: 0, FilePath: "/a"})
	exec.submit(CheckResult{ForecastTime: time.Hour, FilePath: "/b"})
	exec.finish()

	err = exec.wait()

	require.ErrorIs(t, err, ErrCommandTimeout)
}
