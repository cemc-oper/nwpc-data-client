// Package fetcher materializes one resolved data file into a local target
// directory without overwriting an existing destination.
package fetcher

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Method string

const (
	MethodCopy    Method = "copy"
	MethodSymlink Method = "symlink"
	MethodAuto    Method = "auto"
)

var (
	ErrDataNotFound      = errors.New("data not found")
	ErrTargetConflict    = errors.New("target conflict")
	ErrUnsupportedSource = errors.New("unsupported source type")
)

type Config struct {
	TargetDir      string
	Method         Method
	StableInterval time.Duration
}

type Result struct {
	FinalPath     string
	SourcePath    string
	SourceType    string
	LocationLevel string
	Method        Method
	Size          int64
	Reused        bool
}

func (c Config) Validate() error {
	if c.TargetDir == "" {
		return fmt.Errorf("target directory is required")
	}
	if c.StableInterval <= 0 {
		return fmt.Errorf("stable interval must be greater than zero")
	}
	switch c.Method {
	case MethodCopy, MethodSymlink, MethodAuto:
		return nil
	default:
		return fmt.Errorf("method must be copy, symlink, or auto, got %q", c.Method)
	}
}

func hasStorageLevel(level string) bool {
	return slices.Contains(strings.Split(level, "/"), "storage")
}
