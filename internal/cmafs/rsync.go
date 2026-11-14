package cmafs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/fetcher"
)

// Fetch downloads one remote file to a caller-owned temporary local path.
// It never chooses or publishes the final target path.
func (b *Backend) Fetch(ctx context.Context, remotePath string, localTempPath string) error {
	if err := validateRemotePath(remotePath); err != nil {
		return err
	}
	if localTempPath == "" {
		return fmt.Errorf("local temporary path is required")
	}
	commandCtx, cancel := context.WithTimeout(ctx, b.transferTimeout)
	defer cancel()

	remoteSpec := b.destination() + ":" + remotePath
	cmd := exec.CommandContext(commandCtx, b.settings.RsyncPath,
		"--protect-args", "-t", "-e", b.sshCommand(), "--", remoteSpec, localTempPath)
	configureProcessGroup(cmd)
	cmd.WaitDelay = 5 * time.Second
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if errors.Is(commandCtx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("%w: fetch %s timed out after %s", fetcher.ErrBackendOperation, remotePath, b.transferTimeout)
	}
	if err != nil {
		return fmt.Errorf("%w: rsync %s: %v: %s", fetcher.ErrBackendOperation, remotePath, err, outputSummary(stderr.String()))
	}
	return nil
}
