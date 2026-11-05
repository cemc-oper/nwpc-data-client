package checker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func reportTestResult() *RunResult {
	return &RunResult{
		StartTime:    time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC),
		DataType:     "demo/type",
		CheckMode:    CheckModeBatched,
		ArrivalOrder: "unordered",
		StartedAt:    time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC),
		FinishedAt:   time.Date(2026, 10, 26, 0, 1, 0, 0, time.UTC),
		Status:       RunStatusFailed,
		Results: []CheckResult{
			{ForecastTime: 0, Status: CheckStatusSuccess, FilePath: "/data/0", CheckCount: 3, WaitDuration: 1500 * time.Millisecond, CommandDuration: 100 * time.Millisecond},
			{ForecastTime: time.Hour, Status: CheckStatusCheckFailed, Error: errors.New("missing")},
			{ForecastTime: 2 * time.Hour, Status: CheckStatusCommandFailed, Error: errors.New("exit 1")},
			{ForecastTime: 3 * time.Hour, Status: CheckStatusCancelled, Error: context.Canceled},
			{ForecastTime: 4 * time.Hour, Status: CheckStatusPending},
		},
	}
}

func TestWriteRunReportSchemaAndSummary(t *testing.T) {
	reportPath := filepath.Join(t.TempDir(), "report.json")
	require.NoError(t, WriteRunReport(reportPath, reportTestResult()))

	data, err := os.ReadFile(reportPath)
	require.NoError(t, err)
	var report RunReport
	require.NoError(t, json.Unmarshal(data, &report))
	assert.Equal(t, runReportVersion, report.ReportVersion)
	assert.Equal(t, "2026102600", report.StartTime)
	assert.Equal(t, "demo/type", report.DataType)
	assert.Equal(t, CheckModeBatched, report.CheckMode)
	assert.Equal(t, "2026-10-26T00:00:00Z", report.StartedAt)
	assert.Equal(t, "2026-10-26T00:01:00Z", report.FinishedAt)
	assert.Equal(t, ReportSummary{Total: 5, Succeeded: 1, Failed: 2, Cancelled: 1, Pending: 1}, report.Summary)
	require.Len(t, report.Results, 5)
	assert.Equal(t, "000h", report.Results[0].ForecastTime)
	assert.Equal(t, 1.5, report.Results[0].WaitSeconds)
	assert.Equal(t, 0.1, report.Results[0].CommandSeconds)
	assert.Equal(t, "missing", report.Results[1].Error)
	assert.False(t, bytes.Contains(data, []byte("stdout")))
	assert.False(t, bytes.Contains(data, []byte("stderr")))
	assert.False(t, bytes.Contains(data, []byte("environment")))
}

func TestWriteRunReportEmptyPathIsDisabled(t *testing.T) {
	assert.NoError(t, WriteRunReport("", nil))
}

func TestWriteRunReportMissingDirectory(t *testing.T) {
	reportPath := filepath.Join(t.TempDir(), "missing", "report.json")
	err := WriteRunReport(reportPath, reportTestResult())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create temporary run report")
	assert.NoFileExists(t, reportPath)
}

func TestWriteRunReportRenameFailureCleansTemporaryFile(t *testing.T) {
	parent := t.TempDir()
	targetDirectory := filepath.Join(parent, "report.json")
	require.NoError(t, os.Mkdir(targetDirectory, 0755))

	err := WriteRunReport(targetDirectory, reportTestResult())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rename temporary run report")
	matches, globErr := filepath.Glob(filepath.Join(parent, ".report.json.tmp-*"))
	require.NoError(t, globErr)
	assert.Empty(t, matches)
	assert.DirExists(t, targetDirectory)
}

func TestWriteRunReportUnavailableResult(t *testing.T) {
	reportPath := filepath.Join(t.TempDir(), "report.json")
	assert.Error(t, WriteRunReport(reportPath, nil))
	assert.NoFileExists(t, reportPath)
}
