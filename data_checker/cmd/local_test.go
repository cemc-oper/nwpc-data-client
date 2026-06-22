package cmd

import (
	"strings"
	"testing"

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
