package checker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func trackerTestPlan(count int) *runPlan {
	plan := &runPlan{
		config:       CheckerConfig{DataType: "demo/type", CheckMode: CheckModePerTime},
		startTime:    time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC),
		arrivalOrder: "unordered",
	}
	for i := 0; i < count; i++ {
		plan.targets = append(plan.targets, &checkTarget{
			index:        i,
			forecastTime: time.Duration(i) * time.Hour,
			defaultPath:  "NOTFOUND",
		})
	}
	return plan
}

func TestRunTrackerFinalStatesAndSummary(t *testing.T) {
	tracker := newRunTracker(trackerTestPlan(4))
	tracker.start(0)
	tracker.progress(0, "/data/0", 3)
	tracker.completeCheck(CheckResult{targetIndex: 0, FilePath: "/data/0", CheckCount: 3})
	tracker.completeCommand(0, 100*time.Millisecond, nil)

	tracker.start(1)
	tracker.completeCheck(CheckResult{targetIndex: 1, Error: errors.New("missing"), CheckCount: 2})
	tracker.start(2)
	tracker.progress(2, "/data/2", 1)

	result := tracker.finish(errors.New("fail fast"))

	assert.Equal(t, RunStatusFailed, result.Status)
	require.Len(t, result.Results, 4)
	assert.Equal(t, CheckStatusSuccess, result.Results[0].Status)
	assert.Equal(t, 3, result.Results[0].CheckCount)
	assert.Equal(t, 100*time.Millisecond, result.Results[0].CommandDuration)
	assert.Equal(t, CheckStatusCheckFailed, result.Results[1].Status)
	assert.Equal(t, CheckStatusCancelled, result.Results[2].Status)
	assert.Equal(t, CheckStatusPending, result.Results[3].Status)
	assert.Zero(t, result.Results[3].CheckCount)

	report := newRunReport(result)
	assert.Equal(t, ReportSummary{Total: 4, Succeeded: 1, Failed: 1, Cancelled: 1, Pending: 1}, report.Summary)
	for i, item := range report.Results {
		assert.Equal(t, fmt.Sprintf("%03dh", i), item.ForecastTime)
	}
}

func TestRunTrackerParentCancellation(t *testing.T) {
	tracker := newRunTracker(trackerTestPlan(2))
	tracker.start(0)
	tracker.progress(0, "NOTFOUND", 4)

	result := tracker.finish(context.Canceled)

	assert.Equal(t, RunStatusCancelled, result.Status)
	assert.Equal(t, CheckStatusCancelled, result.Results[0].Status)
	assert.ErrorIs(t, result.Results[0].Error, context.Canceled)
	assert.Equal(t, CheckStatusPending, result.Results[1].Status)
}

func TestRunTrackerCommandTimeoutIsFailure(t *testing.T) {
	tracker := newRunTracker(trackerTestPlan(1))
	tracker.start(0)
	tracker.completeCheck(CheckResult{targetIndex: 0})
	timeoutErr := fmt.Errorf("wrapped: %w", ErrCommandTimeout)
	tracker.completeCommand(0, time.Second, timeoutErr)

	result := tracker.finish(timeoutErr)

	assert.Equal(t, RunStatusFailed, result.Status)
	assert.Equal(t, CheckStatusCommandFailed, result.Results[0].Status)
	assert.ErrorIs(t, result.Results[0].Error, ErrCommandTimeout)
}

func TestRunTrackerConcurrentUpdates(t *testing.T) {
	tracker := newRunTracker(trackerTestPlan(16))
	var wg sync.WaitGroup
	for i := range tracker.result.Results {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			tracker.start(index)
			tracker.progress(index, fmt.Sprintf("/data/%d", index), index+1)
			tracker.completeCheck(CheckResult{
				targetIndex: index,
				FilePath:    fmt.Sprintf("/data/%d", index),
				CheckCount:  index + 1,
			})
			tracker.completeCommand(index, time.Millisecond, nil)
		}(i)
	}
	wg.Wait()

	result := tracker.finish(nil)

	assert.Equal(t, RunStatusSuccess, result.Status)
	for i, item := range result.Results {
		assert.Equal(t, CheckStatusSuccess, item.Status)
		assert.Equal(t, i+1, item.CheckCount)
	}
}
