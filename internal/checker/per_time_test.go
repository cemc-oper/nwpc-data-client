package checker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadForecastTimesFromStdinHourLevel(t *testing.T) {
	input := "000h 001h 010h 120h"
	reader := strings.NewReader(input)

	expected := []string{"000h", "001h", "010h", "120h"}

	result, err := readForecastTimesFromStdin(reader)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadForecastTimesFromStdinMinuteLevel(t *testing.T) {
	input := "000h00m 000h10m 001h00m 001h10m"
	reader := strings.NewReader(input)

	expected := []string{"000h00m", "000h10m", "001h00m", "001h10m"}

	result, err := readForecastTimesFromStdin(reader)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadForecastTimesFromStdinEmpty(t *testing.T) {
	reader := strings.NewReader("")

	result, err := readForecastTimesFromStdin(reader)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestRunStartupCancellation(t *testing.T) {
	dataConfigPath := filepath.Join(t.TempDir(), "data.yaml")
	require.NoError(t, os.WriteFile(dataConfigPath, []byte("default: NOTFOUND\npaths: []\n"), 0644))
	config := DefaultCheckerConfig()
	config.DataConfigFile = dataConfigPath
	config.ForecastTimes = []string{"0h"}
	config.CheckInterval = "1h"
	config.CheckIntervalJitter = 0
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Run(ctx, config, time.Now(), strings.NewReader(""))

	assert.ErrorIs(t, err, context.Canceled)
}

func TestCheckDataFileCancellationDuringDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	exec := newCommandExecutor(ctx, 0, 2, nil, time.Now(), "", nil, nil)
	plan := &runPlan{
		config:        DefaultCheckerConfig(),
		startTime:     time.Now(),
		checkInterval: time.Hour,
		delayTime:     time.Hour,
		targets: []*checkTarget{
			{forecastTime: 0, dataConfig: dataconfig.DataConfig{Default: "NOTFOUND"}},
			{forecastTime: time.Hour, dataConfig: dataconfig.DataConfig{Default: "NOTFOUND"}},
		},
	}
	plan.tracker = newRunTracker(plan)
	done := make(chan error, 1)
	go func() {
		err := checkDataFile(exec.context(), plan, exec)
		done <- err
	}()

	// Wait until the first target is polling while the second is still delayed.
	// This avoids racing a zero-interval polling loop against a fixed sleep.
	require.Eventually(t, func() bool {
		plan.tracker.mu.Lock()
		defer plan.tracker.mu.Unlock()
		return plan.tracker.result.Results[0].Status == CheckStatusWaiting &&
			plan.tracker.result.Results[0].CheckCount > 0 &&
			plan.tracker.result.Results[1].Status == CheckStatusPending
	}, time.Second, 5*time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("per-time delay did not return after cancellation")
	}
	assert.False(t, exec.submit(CheckResult{}))
	assert.Equal(t, CheckStatusPending, plan.tracker.result.Results[1].Status)
}

func TestCheckDataFileCancellationDuringPolling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	exec := newCommandExecutor(ctx, 0, 1, nil, time.Now(), "", nil, nil)
	plan := &runPlan{
		config:        DefaultCheckerConfig(),
		startTime:     time.Now(),
		checkInterval: time.Hour,
		targets: []*checkTarget{{
			forecastTime: 0,
			dataConfig:   dataconfig.DataConfig{Default: "NOTFOUND"},
		}},
	}
	done := make(chan error, 1)
	go func() {
		err := checkDataFile(exec.context(), plan, exec)
		done <- err
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("per-time polling did not return after cancellation")
	}
}
