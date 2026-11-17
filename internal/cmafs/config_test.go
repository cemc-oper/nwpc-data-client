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

func TestLoadBackendConfigExpandsHomePaths(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	path := writeBackendConfig(t, `version: 1
cmafs:
  host: 10.40.150.126
  identity_file: ~/.ssh/id_rsa
  known_hosts_file: ~/.ssh/known_hosts
`)

	backend, err := Load(path)
	require.NoError(t, err)
	expectedIdentity := filepath.Join(home, ".ssh", "id_rsa")
	expectedKnownHosts := filepath.Join(home, ".ssh", "known_hosts")
	assert.Equal(t, expectedIdentity, backend.settings.IdentityFile)
	assert.Equal(t, expectedKnownHosts, backend.settings.KnownHostsFile)
	assert.Contains(t, backend.sshOptions(), expectedIdentity)
	assert.Contains(t, backend.sshOptions(), "UserKnownHostsFile="+expectedKnownHosts)
}

func TestExpandHomePathPreservesAbsoluteAndEmptyPaths(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	absolute := filepath.Join(home, "absolute-known-hosts")

	for _, testCase := range []struct {
		name     string
		value    string
		expected string
	}{
		{name: "empty", value: "", expected: ""},
		{name: "absolute", value: absolute, expected: absolute},
		{name: "home", value: "~", expected: filepath.Clean(home)},
		{name: "relative left for validation", value: "relative", expected: "relative"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := expandHomePath("known_hosts_file", testCase.value)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
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
		"other user home":     "version: 1\ncmafs:\n  host: host\n  known_hosts_file: ~other/.ssh/known_hosts\n",
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
