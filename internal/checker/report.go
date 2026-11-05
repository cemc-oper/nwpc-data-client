package checker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
)

const runReportVersion = 1

// RunReport is the versioned JSON document written for a completed run.
type RunReport struct {
	ReportVersion int              `json:"report_version"`
	StartTime     string           `json:"start_time"`
	DataType      string           `json:"data_type"`
	CheckMode     string           `json:"check_mode"`
	ArrivalOrder  string           `json:"arrival_order"`
	StartedAt     string           `json:"started_at"`
	FinishedAt    string           `json:"finished_at"`
	Status        string           `json:"status"`
	Summary       ReportSummary    `json:"summary"`
	Results       []ForecastReport `json:"results"`
}

type ReportSummary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	Cancelled int `json:"cancelled"`
	Pending   int `json:"pending"`
}

type ForecastReport struct {
	ForecastTime   string  `json:"forecast_time"`
	Status         string  `json:"status"`
	FilePath       string  `json:"file_path"`
	CheckCount     int     `json:"check_count"`
	WaitSeconds    float64 `json:"wait_seconds"`
	CommandSeconds float64 `json:"command_seconds"`
	Error          string  `json:"error"`
}

func newRunReport(result *RunResult) RunReport {
	report := RunReport{
		ReportVersion: runReportVersion,
		StartTime:     result.StartTime.Format("2006010215"),
		DataType:      result.DataType,
		CheckMode:     result.CheckMode,
		ArrivalOrder:  result.ArrivalOrder,
		StartedAt:     result.StartedAt.Format(time.RFC3339Nano),
		FinishedAt:    result.FinishedAt.Format(time.RFC3339Nano),
		Status:        result.Status,
		Results:       make([]ForecastReport, 0, len(result.Results)),
	}
	for _, item := range result.Results {
		errorString := ""
		if item.Error != nil {
			errorString = item.Error.Error()
		}
		report.Results = append(report.Results, ForecastReport{
			ForecastTime:   dataconfig.FormatForecastTimeShort(item.ForecastTime),
			Status:         item.Status,
			FilePath:       item.FilePath,
			CheckCount:     item.CheckCount,
			WaitSeconds:    item.WaitDuration.Seconds(),
			CommandSeconds: item.CommandDuration.Seconds(),
			Error:          errorString,
		})
		report.Summary.Total++
		switch item.Status {
		case CheckStatusSuccess:
			report.Summary.Succeeded++
		case CheckStatusCheckFailed, CheckStatusCommandFailed:
			report.Summary.Failed++
		case CheckStatusCancelled:
			report.Summary.Cancelled++
		case CheckStatusPending:
			report.Summary.Pending++
		}
	}
	return report
}

// WriteRunReport atomically writes a complete JSON report in the target
// directory. An empty path disables reporting. No command output or process
// environment is included in the document.
func WriteRunReport(filePath string, result *RunResult) error {
	if filePath == "" {
		return nil
	}
	if result == nil {
		return fmt.Errorf("write run report %s failed: run result is unavailable", filePath)
	}
	data, err := json.MarshalIndent(newRunReport(result), "", "  ")
	if err != nil {
		return fmt.Errorf("encode run report %s failed: %w", filePath, err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(filePath)
	tempFile, err := os.CreateTemp(dir, "."+filepath.Base(filePath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temporary run report for %s failed: %w", filePath, err)
	}
	tempPath := tempFile.Name()
	committed := false
	defer func() {
		_ = tempFile.Close()
		if !committed {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("write temporary run report for %s failed: %w", filePath, err)
	}
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("sync temporary run report for %s failed: %w", filePath, err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary run report for %s failed: %w", filePath, err)
	}
	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("rename temporary run report to %s failed: %w", filePath, err)
	}
	committed = true
	return nil
}
