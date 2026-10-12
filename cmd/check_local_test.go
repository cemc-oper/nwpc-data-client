package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadForecastTimesFromStdinHourLevel(t *testing.T) {
	input := "000h 001h 010h 120h"
	reader := strings.NewReader(input)

	expected := []string{"000h", "001h", "010h", "120h"}

	result := readForecastTimesFromStdin(reader)

	assert.Equal(t, expected, result)
}

func TestReadForecastTimesFromStdinMinuteLevel(t *testing.T) {
	input := "000h00m 000h10m 001h00m 001h10m"
	reader := strings.NewReader(input)

	expected := []string{"000h00m", "000h10m", "001h00m", "001h10m"}

	result := readForecastTimesFromStdin(reader)

	assert.Equal(t, expected, result)
}

func TestReadForecastTimesFromStdinEmpty(t *testing.T) {
	reader := strings.NewReader("")

	result := readForecastTimesFromStdin(reader)

	assert.Empty(t, result)
}

func TestWaitForStableFileAlreadyStable(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	content := []byte("some data content")
	err := os.WriteFile(filePath, content, 0644)
	assert.NoError(t, err)

	size, err := waitForStableFile(filePath, 10*time.Millisecond, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(content)), size)
}

func TestWaitForStableFileStableAfterGrowth(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	err := os.WriteFile(filePath, []byte("chunk0"), 0644)
	assert.NoError(t, err)

	// Grow the file during the first checks, then leave it alone.
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 4; i++ {
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			_, _ = f.WriteString("chunk")
			_ = f.Close()
			time.Sleep(10 * time.Millisecond)
		}
		<-stop
	}()

	size, err := waitForStableFile(filePath, 10*time.Millisecond, 50)
	close(stop)
	<-done

	assert.NoError(t, err)
	info, statErr := os.Stat(filePath)
	assert.NoError(t, statErr)
	assert.Equal(t, info.Size(), size)
	assert.Equal(t, int64(len("chunk0")+4*len("chunk")), size)
}

func TestWaitForStableFileKeepsGrowing(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	err := os.WriteFile(filePath, []byte("chunk0"), 0644)
	assert.NoError(t, err)

	// Keep growing the file so it never stabilizes within the check budget.
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
			}
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			_, _ = f.WriteString("chunk")
			_ = f.Close()
			time.Sleep(5 * time.Millisecond)
		}
	}()

	_, err = waitForStableFile(filePath, 10*time.Millisecond, 5)
	close(stop)
	<-done

	assert.Error(t, err)
}

func TestWaitForStableFileMissing(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-exist.grb2")

	_, err := waitForStableFile(filePath, 10*time.Millisecond, 5)

	assert.Error(t, err)
}
