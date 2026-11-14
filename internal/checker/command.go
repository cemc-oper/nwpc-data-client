package checker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/commander-cli/cmd"
	log "github.com/sirupsen/logrus"
)

// ErrCommandTimeout identifies a deadline created for an individual command,
// as distinct from cancellation or a deadline on the parent run context.
var ErrCommandTimeout = errors.New("command timed out")

type CheckerTemplateVariable struct {
	dataconfig.ConfigTemplateVariable
	FilePath string
}

func runCommand(ctx context.Context, timeout time.Duration, commandTemplate *template.Template, startTime time.Time, forecastTime time.Duration, member string, filePath string) error {
	return runCommandWithStorageBase(ctx, timeout, commandTemplate, startTime, forecastTime, member, "", filePath)
}

func runCommandWithStorageBase(ctx context.Context, timeout time.Duration, commandTemplate *template.Template, startTime time.Time, forecastTime time.Duration, member string, storageBase string, filePath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tpVar := dataconfig.GenerateConfigTemplateVariableWithStorageBase(startTime, forecastTime, member, storageBase)
	var checkerVar CheckerTemplateVariable
	checkerVar.ConfigTemplateVariable = tpVar
	checkerVar.FilePath = filePath

	var commandBuilder strings.Builder
	err := commandTemplate.Execute(&commandBuilder, checkerVar)
	if err != nil {
		return fmt.Errorf("command template execute has error: %v", err)
	}

	commandString := commandBuilder.String()

	forecastTimeString := dataconfig.FormatForecastTimeShort(forecastTime)
	log.WithFields(log.Fields{
		"forecast_time": forecastTimeString,
	}).Infof("running command <%s> ...", commandString)

	commandCtx := ctx
	cancel := func() {}
	if timeout > 0 {
		commandCtx, cancel = context.WithTimeout(ctx, timeout)
	}
	defer cancel()

	c := cmd.NewCommand(
		commandString,
		cmd.WithStandardStreams,
		cmd.WithInheritedEnvironment(nil),
		// The dependency's implicit 30-minute timer is replaced by the
		// explicit per-command context so command_timeout=0 really disables it.
		cmd.WithoutTimeout)
	err = c.ExecuteContext(commandCtx)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("run command <%s> cancelled: %w", commandString, ctx.Err())
		}
		if errors.Is(commandCtx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("%w after %s: <%s>: %v", ErrCommandTimeout, timeout, commandString, err)
		}
		return fmt.Errorf("run command <%s> has error: %w", commandString, err)
	} else if c.ExitCode() != 0 {
		return fmt.Errorf("run command <%s> exit code is not 0: %d", commandString, c.ExitCode())
	}

	return nil
}

func runCommands(ctx context.Context, timeout time.Duration, commandTemplates []*template.Template, startTime time.Time, forecastTime time.Duration, member string, storageBase string, filePath string) error {
	for _, commandTemplate := range commandTemplates {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := runCommandWithStorageBase(ctx, timeout, commandTemplate, startTime, forecastTime, member, storageBase, filePath); err != nil {
			return err
		}
	}
	return nil
}

func buildCommandTemplates(command string, commands []string) ([]*template.Template, error) {
	if command != "" && len(commands) > 0 {
		return nil, fmt.Errorf("execute-command and execute-commands are mutually exclusive")
	}

	funcMap := template.FuncMap{
		"generateStartTime":    dataconfig.GenerateStartTime,
		"getYear":              dataconfig.GetYear,
		"getMonth":             dataconfig.GetMonth,
		"getDay":               dataconfig.GetDay,
		"getHour":              dataconfig.GetHour,
		"generateForecastTime": dataconfig.GenerateForecastTime,
		"getForecastHour":      dataconfig.GetForecastHour,
		"getForecastMinute":    dataconfig.GetForecastMinute,
	}

	var sources []string
	if command != "" {
		sources = []string{command}
	} else {
		sources = commands
	}

	var templates []*template.Template
	for i, s := range sources {
		t, err := template.New(fmt.Sprintf("command-%d", i)).Funcs(funcMap).Parse(s)
		if err != nil {
			return nil, fmt.Errorf("parse execute command failed: %v", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}
