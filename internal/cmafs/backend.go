package cmafs

import (
	"fmt"
	"path"
	"strings"
	"unicode"
)

func validateRemotePath(remotePath string) error {
	if !path.IsAbs(remotePath) {
		return fmt.Errorf("cmafs path must be absolute: %q", remotePath)
	}
	for _, r := range remotePath {
		if unicode.IsControl(r) {
			return fmt.Errorf("cmafs path must not contain control characters")
		}
	}
	return nil
}

func (b *Backend) destination() string {
	host := b.settings.Host
	if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
		host = "[" + host + "]"
	}
	if b.settings.User != "" {
		return b.settings.User + "@" + host
	}
	return host
}

// sshOptions deliberately omits -n: rsync needs SSH stdin/stdout for its
// protocol. Standalone SSH Stat adds -n separately.
func (b *Backend) sshOptions() []string {
	connectSeconds := max(1, int(b.connectTimeout.Seconds()))
	args := []string{
		"-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=yes",
		"-o", fmt.Sprintf("ConnectTimeout=%d", connectSeconds),
		"-o", "ServerAliveInterval=30",
		"-o", "ServerAliveCountMax=3",
		"-p", fmt.Sprintf("%d", b.settings.Port),
	}
	if b.settings.KnownHostsFile != "" {
		args = append(args, "-o", "UserKnownHostsFile="+b.settings.KnownHostsFile)
	}
	if b.settings.IdentityFile != "" {
		args = append(args, "-i", b.settings.IdentityFile)
	}
	return args
}

func (b *Backend) sshCommand() string {
	parts := []string{shellQuote(b.settings.SSHPath)}
	for _, arg := range b.sshOptions() {
		parts = append(parts, shellQuote(arg))
	}
	return strings.Join(parts, " ")
}

// shellQuote protects an argument interpreted by a POSIX shell, including the
// remote command passed through SSH and the rsync -e transport string.
func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func outputSummary(value string) string {
	const maxLength = 4096
	value = strings.TrimSpace(value)
	if len(value) > maxLength {
		return value[:maxLength] + "..."
	}
	return value
}
