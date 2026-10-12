package checker

import (
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
// Fail-fast semantics match the historical per-time consumer: a check error
// or a failing command terminates the process via log.Fatal. The onTaskDone
// hook (nil for per-time mode) runs after a result has been fully processed
// and is the attachment point for sequential arrival gating.
type commandExecutor struct {
	queue            chan CheckResult
	drained          chan struct{}
	commandTemplates []*template.Template
	startTime        time.Time
	member           string
	onTaskDone       func(CheckResult)
}

func newCommandExecutor(capacity int, commandTemplates []*template.Template, startTime time.Time, member string, onTaskDone func(CheckResult)) *commandExecutor {
	e := &commandExecutor{
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

// submit enqueues a check result. It never blocks as long as the queue
// capacity covers the total number of forecast times.
func (e *commandExecutor) submit(result CheckResult) {
	e.queue <- result
}

// finish closes the queue; no more submissions are allowed afterwards.
func (e *commandExecutor) finish() {
	close(e.queue)
}

// wait blocks until the queue is closed and fully drained.
func (e *commandExecutor) wait() {
	<-e.drained
}

func (e *commandExecutor) consume() {
	defer close(e.drained)
	for result := range e.queue {
		forecastTimeString := dataconfig.FormatForecastTimeShort(result.ForecastTime)
		currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
		if result.Error != nil {
			currentLog.Fatalf("check failed: %v", result.Error)
		}

		currentLog.Infof("file is available, run command...")

		if len(e.commandTemplates) == 0 {
			currentLog.Debugf("ignore execute command because of empty command")
		} else {
			err := runCommands(e.commandTemplates, e.startTime, result.ForecastTime, e.member, result.FilePath)
			if err != nil {
				currentLog.Fatalf("run command failed: %v", err)
			}
			currentLog.Infof("run command success")
		}

		if e.onTaskDone != nil {
			e.onTaskDone(result)
		}
	}
}
