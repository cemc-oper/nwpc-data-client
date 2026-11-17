// Package cmafs implements read-only secondary-storage access through the
// same SSH and rsync transport used by the cmafs reference tools.
package cmafs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

const CurrentBackendConfigVersion = 1

type fileConfig struct {
	Version int       `yaml:"version"`
	Cmafs   *Settings `yaml:"cmafs"`
}

// Settings contains machine-level connection parameters. Passwords and private
// key contents intentionally have no representation in this configuration.
type Settings struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	SSHPath         string `yaml:"ssh_path"`
	RsyncPath       string `yaml:"rsync_path"`
	IdentityFile    string `yaml:"identity_file"`
	KnownHostsFile  string `yaml:"known_hosts_file"`
	ConnectTimeout  string `yaml:"connect_timeout"`
	TransferTimeout string `yaml:"transfer_timeout"`
}

// Backend is a read-only SSH/rsync secondary-storage backend.
type Backend struct {
	settings        Settings
	connectTimeout  time.Duration
	transferTimeout time.Duration
}

// Load parses and validates a backend YAML file without accessing the network.
// A file without a cmafs section returns (nil, nil), meaning no remote backend.
func Load(filePath string) (*Backend, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read backend config %s: %w", filePath, err)
	}
	var config fileConfig
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("parse backend config %s: %w", filePath, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, fmt.Errorf("parse backend config %s: multiple YAML documents are not supported", filePath)
		}
		return nil, fmt.Errorf("parse backend config %s: %w", filePath, err)
	}
	if config.Version != CurrentBackendConfigVersion {
		return nil, fmt.Errorf("backend config %s: version must be %d", filePath, CurrentBackendConfigVersion)
	}
	if config.Cmafs == nil {
		return nil, nil
	}

	settings := *config.Cmafs
	setDefaults(&settings)
	if err := expandUserPaths(&settings); err != nil {
		return nil, fmt.Errorf("backend config %s: cmafs: %w", filePath, err)
	}
	if err := validateSettings(settings); err != nil {
		return nil, fmt.Errorf("backend config %s: cmafs: %w", filePath, err)
	}
	connectTimeout, err := time.ParseDuration(settings.ConnectTimeout)
	if err != nil || connectTimeout <= 0 {
		return nil, fmt.Errorf("backend config %s: cmafs.connect_timeout must be a positive duration", filePath)
	}
	transferTimeout, err := time.ParseDuration(settings.TransferTimeout)
	if err != nil || transferTimeout <= 0 {
		return nil, fmt.Errorf("backend config %s: cmafs.transfer_timeout must be a positive duration", filePath)
	}
	return &Backend{
		settings:        settings,
		connectTimeout:  connectTimeout,
		transferTimeout: transferTimeout,
	}, nil
}

func setDefaults(settings *Settings) {
	if settings.Port == 0 {
		settings.Port = 22
	}
	if settings.SSHPath == "" {
		settings.SSHPath = "ssh"
	}
	if settings.RsyncPath == "" {
		settings.RsyncPath = "rsync"
	}
	if settings.ConnectTimeout == "" {
		settings.ConnectTimeout = "15s"
	}
	if settings.TransferTimeout == "" {
		settings.TransferTimeout = "2h"
	}
}

// expandUserPaths resolves home-relative file settings for the current user.
// Environment interpolation and other users' homes are deliberately unsupported.
func expandUserPaths(settings *Settings) error {
	fields := []struct {
		name  string
		value *string
	}{
		{name: "identity_file", value: &settings.IdentityFile},
		{name: "known_hosts_file", value: &settings.KnownHostsFile},
	}
	for _, field := range fields {
		expanded, err := expandHomePath(field.name, *field.value)
		if err != nil {
			return err
		}
		*field.value = expanded
	}
	return nil
}

func expandHomePath(name string, value string) (string, error) {
	if value == "" || filepath.IsAbs(value) {
		return value, nil
	}

	homeRelative := value == "~" || strings.HasPrefix(value, "~/")
	if filepath.Separator == '\\' {
		homeRelative = homeRelative || strings.HasPrefix(value, `~\`)
	}
	if !homeRelative {
		if strings.HasPrefix(value, "~") {
			return "", fmt.Errorf("%s must be absolute or start with ~/; other users' home directories are not supported", name)
		}
		return value, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve current user home for %s: %w", name, err)
	}
	if value == "~" {
		return filepath.Clean(home), nil
	}
	return filepath.Join(home, value[2:]), nil
}

func validateSettings(settings Settings) error {
	if err := validateToken("host", settings.Host, false); err != nil {
		return err
	}
	if settings.Port < 1 || settings.Port > 65535 {
		return fmt.Errorf("port must be in [1, 65535]")
	}
	for name, value := range map[string]string{
		"user":       settings.User,
		"ssh_path":   settings.SSHPath,
		"rsync_path": settings.RsyncPath,
	} {
		if err := validateToken(name, value, name == "user"); err != nil {
			return err
		}
	}
	for name, value := range map[string]string{
		"identity_file":    settings.IdentityFile,
		"known_hosts_file": settings.KnownHostsFile,
	} {
		if value != "" && !filepath.IsAbs(value) {
			return fmt.Errorf("%s must be an absolute path", name)
		}
		if strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("%s must not contain newlines", name)
		}
	}
	return nil
}

func validateToken(name string, value string, allowEmpty bool) error {
	if value == "" {
		if allowEmpty {
			return nil
		}
		return fmt.Errorf("%s is required", name)
	}
	if strings.HasPrefix(value, "-") {
		return fmt.Errorf("%s must not start with '-'", name)
	}
	for _, r := range value {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return fmt.Errorf("%s must not contain whitespace or control characters", name)
		}
	}
	return nil
}
