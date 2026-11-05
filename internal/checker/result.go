package checker

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	RunStatusSuccess   = "success"
	RunStatusFailed    = "failed"
	RunStatusCancelled = "cancelled"

	CheckStatusPending       = "pending"
	CheckStatusWaiting       = "waiting"
	CheckStatusSuccess       = "success"
	CheckStatusCheckFailed   = "check_failed"
	CheckStatusCommandFailed = "command_failed"
	CheckStatusCancelled     = "cancelled"
)

// CheckResult is the final state of one forecast time. The unexported fields
// connect queue entries and timing state to a stable plan position.
type CheckResult struct {
	ForecastTime    time.Duration
	FilePath        string
	Error           error
	Status          string
	CheckCount      int
	WaitDuration    time.Duration
	CommandDuration time.Duration

	targetIndex int
	waitStarted time.Time
	waitEnded   time.Time
}

// RunResult is a final, internally consistent snapshot of a checker run.
type RunResult struct {
	StartTime    time.Time
	DataType     string
	CheckMode    string
	ArrivalOrder string
	StartedAt    time.Time
	FinishedAt   time.Time
	Status       string
	Results      []CheckResult
	Err          error
}

type runTracker struct {
	mu     sync.Mutex
	result RunResult
}

func newRunTracker(plan *runPlan) *runTracker {
	startedAt := time.Now()
	results := make([]CheckResult, len(plan.targets))
	for i, target := range plan.targets {
		results[i] = CheckResult{
			ForecastTime: target.forecastTime,
			FilePath:     target.defaultPath,
			Status:       CheckStatusPending,
			targetIndex:  i,
		}
	}
	return &runTracker{result: RunResult{
		StartTime:    plan.startTime,
		DataType:     plan.config.DataType,
		CheckMode:    plan.config.CheckMode,
		ArrivalOrder: plan.arrivalOrder,
		StartedAt:    startedAt.UTC(),
		Status:       CheckStatusPending,
		Results:      results,
	}}
}

func (t *runTracker) start(index int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if index < 0 || index >= len(t.result.Results) {
		return
	}
	result := &t.result.Results[index]
	if result.Status != CheckStatusPending {
		return
	}
	result.Status = CheckStatusWaiting
	result.waitStarted = time.Now()
}

func (t *runTracker) progress(index int, filePath string, checkCount int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if index < 0 || index >= len(t.result.Results) {
		return
	}
	result := &t.result.Results[index]
	result.FilePath = filePath
	result.CheckCount = checkCount
	if !result.waitStarted.IsZero() && result.waitEnded.IsZero() {
		result.WaitDuration = time.Since(result.waitStarted)
	}
}

func (t *runTracker) completeCheck(result CheckResult) {
	t.mu.Lock()
	defer t.mu.Unlock()
	index := result.targetIndex
	if index < 0 || index >= len(t.result.Results) {
		return
	}
	current := &t.result.Results[index]
	current.FilePath = result.FilePath
	current.CheckCount = result.CheckCount
	current.Error = result.Error
	if current.waitStarted.IsZero() {
		current.waitStarted = time.Now()
	}
	if current.waitEnded.IsZero() {
		current.waitEnded = time.Now()
		current.WaitDuration = current.waitEnded.Sub(current.waitStarted)
	}
	if result.Error != nil {
		current.Status = CheckStatusCheckFailed
	} else {
		// The file is ready; it remains waiting until the serial executor has
		// processed its command or the run is cancelled.
		current.Status = CheckStatusWaiting
	}
}

func (t *runTracker) completeCommand(index int, duration time.Duration, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if index < 0 || index >= len(t.result.Results) {
		return
	}
	result := &t.result.Results[index]
	result.CommandDuration = duration
	result.Error = err
	switch {
	case err == nil:
		result.Status = CheckStatusSuccess
	case isRunCancellation(err):
		result.Status = CheckStatusCancelled
	default:
		result.Status = CheckStatusCommandFailed
	}
}

func isRunCancellation(err error) bool {
	return err != nil && !errors.Is(err, ErrCommandTimeout) &&
		(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded))
}

func (t *runTracker) finish(runErr error) *RunResult {
	t.mu.Lock()
	defer t.mu.Unlock()
	finishedAt := time.Now()
	for i := range t.result.Results {
		result := &t.result.Results[i]
		if result.Status != CheckStatusWaiting {
			continue
		}
		if result.waitEnded.IsZero() && !result.waitStarted.IsZero() {
			result.WaitDuration = finishedAt.Sub(result.waitStarted)
		}
		result.Status = CheckStatusCancelled
		if isRunCancellation(runErr) {
			result.Error = runErr
		} else {
			result.Error = context.Canceled
		}
	}

	t.result.FinishedAt = finishedAt.UTC()
	t.result.Err = runErr
	switch {
	case runErr == nil:
		t.result.Status = RunStatusSuccess
	case isRunCancellation(runErr):
		t.result.Status = RunStatusCancelled
	default:
		t.result.Status = RunStatusFailed
	}
	result := t.result
	result.Results = append([]CheckResult(nil), t.result.Results...)
	return &result
}
