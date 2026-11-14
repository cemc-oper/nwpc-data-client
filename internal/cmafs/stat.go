package cmafs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/fetcher"
)

// Stat returns remote type, size, and mtime using a fixed GNU stat protocol.
// Only explicit ENOENT/ENOTDIR text maps to fetcher.ErrSourceNotExist.
func (b *Backend) Stat(ctx context.Context, remotePath string) (fetcher.RemoteMetadata, error) {
	if err := validateRemotePath(remotePath); err != nil {
		return fetcher.RemoteMetadata{}, err
	}
	commandCtx, cancel := context.WithTimeout(ctx, b.connectTimeout)
	defer cancel()

	args := []string{"-n"}
	args = append(args, b.sshOptions()...)
	args = append(args, "--", b.destination(),
		"LC_ALL=C stat -L -c '%F|%s|%Y' -- "+shellQuote(remotePath))
	cmd := exec.CommandContext(commandCtx, b.settings.SSHPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		return fetcher.RemoteMetadata{}, ctx.Err()
	}
	if errors.Is(commandCtx.Err(), context.DeadlineExceeded) {
		return fetcher.RemoteMetadata{}, fmt.Errorf("%w: stat %s timed out after %s", fetcher.ErrBackendOperation, remotePath, b.connectTimeout)
	}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			message := outputSummary(stderr.String())
			if strings.Contains(message, "No such file or directory") || strings.Contains(message, "Not a directory") {
				return fetcher.RemoteMetadata{}, fmt.Errorf("%w: %s", fetcher.ErrSourceNotExist, remotePath)
			}
		}
		return fetcher.RemoteMetadata{}, fmt.Errorf("%w: stat %s: %v: %s", fetcher.ErrBackendOperation, remotePath, err, outputSummary(stderr.String()))
	}

	line := strings.TrimSpace(stdout.String())
	fields := strings.Split(line, "|")
	if len(fields) != 3 || strings.Contains(line, "\n") {
		return fetcher.RemoteMetadata{}, fmt.Errorf("%w: stat %s returned invalid record %q", fetcher.ErrBackendOperation, remotePath, line)
	}
	size, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return fetcher.RemoteMetadata{}, fmt.Errorf("%w: stat %s returned invalid size %q", fetcher.ErrBackendOperation, remotePath, fields[1])
	}
	mtime, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return fetcher.RemoteMetadata{}, fmt.Errorf("%w: stat %s returned invalid mtime %q", fetcher.ErrBackendOperation, remotePath, fields[2])
	}
	return fetcher.RemoteMetadata{
		Size:    size,
		ModTime: time.Unix(mtime, 0),
		Regular: fields[0] == "regular file",
	}, nil
}
