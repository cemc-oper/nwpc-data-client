package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseForecastTimeInputHourLevel(t *testing.T) {
	input := "000h 001h 010h 120h"
	reader := strings.NewReader(input)

	expected := []time.Duration{
		0 * time.Hour,
		1 * time.Hour,
		10 * time.Hour,
		120 * time.Hour,
	}

	result := parseForecastTimeInput(reader)

	assert.Equal(t, expected, result)
}

func TestParseForecastTimeInputMinuteLevel(t *testing.T) {
	input := "000h00m 000h10m 001h00m 001h10m"
	reader := strings.NewReader(input)

	expected := []time.Duration{
		0 * time.Hour,
		0*time.Hour + 10*time.Minute,
		1 * time.Hour,
		1*time.Hour + 10*time.Minute,
	}

	result := parseForecastTimeInput(reader)

	assert.Equal(t, expected, result)
}
