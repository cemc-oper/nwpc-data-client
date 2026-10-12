package checker

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/commander-cli/cmd"
	log "github.com/sirupsen/logrus"
)

type CheckerTemplateVariable struct {
	dataconfig.ConfigTemplateVariable
	FilePath string
}

func runCommand(commandTemplate *template.Template, startTime time.Time, forecastTime time.Duration, member string, filePath string) error {
	tpVar := dataconfig.GenerateConfigTemplateVariable(startTime, forecastTime, member)
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

	c := cmd.NewCommand(
		commandString,
		cmd.WithStandardStreams,
		cmd.WithInheritedEnvironment(nil))
	err = c.Execute()
	if err != nil {
		return fmt.Errorf("run command <%s> has error: %v", commandString, err)
	} else if c.ExitCode() != 0 {
		return fmt.Errorf("run command <%s> exit code is not 0: %d", commandString, c.ExitCode())
	}

	return nil
}

func runCommands(commandTemplates []*template.Template, startTime time.Time, forecastTime time.Duration, member string, filePath string) error {
	for _, commandTemplate := range commandTemplates {
		if err := runCommand(commandTemplate, startTime, forecastTime, member, filePath); err != nil {
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
