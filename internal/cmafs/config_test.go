package cmafs

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeBackendConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "backend.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestLoadBackendConfigDefaults(t *testing.T) {
	path := writeBackendConfig(t, "version: 1\ncmafs:\n  host: 10.40.150.126\n")
	backend, err := Load(path)
	require.NoError(t, err)
	require.NotNil(t, backend)
	assert.Equal(t, 22, backend.settings.Port)
	assert.Equal(t, "ssh", backend.settings.SSHPath)
	assert.Equal(t, "rsync", backend.settings.RsyncPath)
	assert.Equal(t, 15*time.Second, backend.connectTimeout)
	assert.Equal(t, 2*time.Hour, backend.transferTimeout)
	assert.Equal(t, "10.40.150.126", backend.destination())
}

func TestLoadBackendConfigWithoutCmafs(t *testing.T) {
	backend, err := Load(writeBackendConfig(t, "version: 1\n"))
	require.NoError(t, err)
	assert.Nil(t, backend)
}

func TestLoadBackendConfigRejectsInvalidValues(t *testing.T) {
	tests := map[string]string{
		"missing version":     "cmafs:\n  host: host\n",
		"unsupported version": "version: 2\ncmafs:\n  host: host\n",
		"missing host":        "version: 1\ncmafs: {}\n",
		"unknown field":       "version: 1\ncmafs:\n  host: host\n  password: secret\n",
		"invalid port":        "version: 1\ncmafs:\n  host: host\n  port: -1\n",
		"relative identity":   "version: 1\ncmafs:\n  host: host\n  identity_file: key\n",
		"invalid connect":     "version: 1\ncmafs:\n  host: host\n  connect_timeout: 0s\n",
		"invalid transfer":    "version: 1\ncmafs:\n  host: host\n  transfer_timeout: nope\n",
		"option host":         "version: 1\ncmafs:\n  host: -oProxyCommand=x\n",
		"space in host":       "version: 1\ncmafs:\n  host: bad host\n",
		"multiple documents":  "version: 1\n---\nversion: 1\n",
	}
	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeBackendConfig(t, content))
			assert.Error(t, err)
		})
	}
}

func TestLoadBackendConfigMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	assert.Error(t, err)
}

func TestShellQuoteAndIPv6Destination(t *testing.T) {
	assert.Equal(t, `'a'"'"'b'`, shellQuote("a'b"))
	backend := &Backend{settings: Settings{Host: "2001:db8::1", User: "user"}}
	assert.Equal(t, "user@[2001:db8::1]", backend.destination())
}
