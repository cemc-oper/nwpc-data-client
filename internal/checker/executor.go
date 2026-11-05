package checker

import (
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
)

// RunResult contains the check results processed by the executor, in queue
// order. Later reporting stages can extend it without changing polling logic.
type RunResult struct {
	Results []CheckResult
}

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
	queue            chan CheckResult
	drained          chan struct{}
	commandTemplates []*template.Template
	startTime        time.Time
	member           string
	onTaskDone       func(CheckResult)
	results          []CheckResult
	err              error
}

func newCommandExecutor(ctx context.Context, capacity int, commandTemplates []*template.Template, startTime time.Time, member string, onTaskDone func(CheckResult)) *commandExecutor {
	runCtx, cancel := context.WithCancel(ctx)
	e := &commandExecutor{
		ctx:              runCtx,
		cancel:           cancel,
		queue:            make(chan CheckResult, capacity),
		drained:          make(chan struct{}),
		commandTemplates: commandTemplates,
		startTime:        startTime,
		member:           member,
		onTaskDone:       onTaskDone,
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

// wait blocks until the consumer exits and returns its processed results and
// first error. The queue must be closed after all producers have returned.
func (e *commandExecutor) wait() (RunResult, error) {
	<-e.drained
	e.cancel()
	return RunResult{Results: e.results}, e.err
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
				result.Error = fmt.Errorf("check failed for forecast time %s: %w", forecastTimeString, result.Error)
				e.results = append(e.results, result)
				e.fail(result.Error)
				return
			}

			currentLog.Infof("file is available, run command...")

			if len(e.commandTemplates) == 0 {
				currentLog.Debugf("ignore execute command because of empty command")
			} else {
				err := runCommands(e.commandTemplates, e.startTime, result.ForecastTime, e.member, result.FilePath)
				if err != nil {
					result.Error = fmt.Errorf("run command failed for forecast time %s: %w", forecastTimeString, err)
					e.results = append(e.results, result)
					e.fail(result.Error)
					return
				}
				currentLog.Infof("run command success")
			}

			e.results = append(e.results, result)
			if e.onTaskDone != nil {
				e.onTaskDone(result)
			}
		}
	}
}
