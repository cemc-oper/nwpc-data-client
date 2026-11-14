package checker

import (
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
)

// commandExecutor runs the configured command(s) for each submitted check
// result, serially and in submission order, from a single consumer goroutine.
// The queue is buffered with capacity for every forecast time, so submit
// never blocks a producer.
//
// A check error, command error, or context cancellation stops processing and
// cancels the run. Producers observe that same context and finish before the
// caller closes the queue. The onTaskDone hook (nil for per-time mode) runs
// after a result has been fully processed and is the attachment point for
// sequential arrival gating.
type commandExecutor struct {
	ctx              context.Context
	cancel           context.CancelFunc
	commandTimeout   time.Duration
	queue            chan CheckResult
	drained          chan struct{}
	commandTemplates []*template.Template
	startTime        time.Time
	member           string
	storageBase      string
	onTaskDone       func(CheckResult)
	tracker          *runTracker
	err              error
}

func newCommandExecutor(ctx context.Context, commandTimeout time.Duration, capacity int, commandTemplates []*template.Template, startTime time.Time, member string, onTaskDone func(CheckResult), tracker *runTracker) *commandExecutor {
	return newCommandExecutorWithStorageBase(ctx, commandTimeout, capacity, commandTemplates, startTime, member, "", onTaskDone, tracker)
}

func newCommandExecutorWithStorageBase(ctx context.Context, commandTimeout time.Duration, capacity int, commandTemplates []*template.Template, startTime time.Time, member string, storageBase string, onTaskDone func(CheckResult), tracker *runTracker) *commandExecutor {
	runCtx, cancel := context.WithCancel(ctx)
	e := &commandExecutor{
		ctx:              runCtx,
		cancel:           cancel,
		commandTimeout:   commandTimeout,
		queue:            make(chan CheckResult, capacity),
		drained:          make(chan struct{}),
		commandTemplates: commandTemplates,
		startTime:        startTime,
		member:           member,
		storageBase:      storageBase,
		onTaskDone:       onTaskDone,
		tracker:          tracker,
	}
	go e.consume()
	return e
}

func (e *commandExecutor) context() context.Context {
	return e.ctx
}

// submit enqueues a check result unless the run has already stopped. It never
// blocks as long as the queue capacity covers the total number of forecast times.
func (e *commandExecutor) submit(result CheckResult) bool {
	if e.ctx.Err() != nil {
		return false
	}
	select {
	case e.queue <- result:
		return true
	case <-e.ctx.Done():
		return false
	}
}

// setOnTaskDone installs the task-done hook. It must be called before any
// result is submitted; it exists so the polling loop can decide on sequential
// gating only after the data config has been expanded.
func (e *commandExecutor) setOnTaskDone(hook func(CheckResult)) {
	e.onTaskDone = hook
}

// finish closes the queue; no more submissions are allowed afterwards.
func (e *commandExecutor) finish() {
	close(e.queue)
}

// wait blocks until the consumer exits and returns its first error. The queue
// must be closed after all producers have returned.
func (e *commandExecutor) wait() error {
	<-e.drained
	e.cancel()
	return e.err
}

func (e *commandExecutor) fail(err error) {
	if e.err == nil {
		e.err = err
	}
	e.cancel()
}

func (e *commandExecutor) consume() {
	defer close(e.drained)
	for {
		if err := e.ctx.Err(); err != nil {
			e.fail(fmt.Errorf("check run cancelled: %w", err))
			return
		}
		select {
		case <-e.ctx.Done():
			e.fail(fmt.Errorf("check run cancelled: %w", e.ctx.Err()))
			return
		case result, ok := <-e.queue:
			if !ok {
				return
			}
			// Do not start another result if cancellation raced with receive.
			if err := e.ctx.Err(); err != nil {
				e.fail(fmt.Errorf("check run cancelled: %w", err))
				return
			}

			forecastTimeString := dataconfig.FormatForecastTimeShort(result.ForecastTime)
			currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
			if result.Error != nil {
				e.fail(fmt.Errorf("check failed for forecast time %s: %w", forecastTimeString, result.Error))
				return
			}

			currentLog.Infof("file is available, run command...")

			if len(e.commandTemplates) == 0 {
				currentLog.Debugf("ignore execute command because of empty command")
				if e.tracker != nil {
					e.tracker.completeCommand(result.targetIndex, 0, nil)
				}
			} else {
				commandStarted := time.Now()
				err := runCommands(e.ctx, e.commandTimeout, e.commandTemplates, e.startTime, result.ForecastTime, e.member, e.storageBase, result.FilePath)
				if e.tracker != nil {
					e.tracker.completeCommand(result.targetIndex, time.Since(commandStarted), err)
				}
				if err != nil {
					e.fail(fmt.Errorf("run command failed for forecast time %s: %w", forecastTimeString, err))
					return
				}
				currentLog.Infof("run command success")
			}

			if e.onTaskDone != nil {
				e.onTaskDone(result)
			}
		}
	}
}
