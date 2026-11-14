//go:build !windows

package cmafs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/fetcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const helperDataSize = 47
const helperMtime = 1653633373

func helperBackend(t *testing.T, mode string) *Backend {
	t.Helper()
	helper := filepath.Join(t.TempDir(), "cmafs-helper")
	script := "#!/bin/sh\nexec " + shellQuote(os.Args[0]) + " -test.run=TestCmafsHelperProcess -- \"$@\"\n"
	require.NoError(t, os.WriteFile(helper, []byte(script), 0755))
	t.Setenv("CMAFS_HELPER_MODE", mode)
	t.Setenv("CMAFS_HELPER_ARGS", filepath.Join(t.TempDir(), "args"))
	t.Setenv("CMAFS_HELPER_CHILD_PID", filepath.Join(t.TempDir(), "child-pid"))
	return &Backend{
		settings: Settings{
			Host:           "server",
			SSHPath:        helper,
			RsyncPath:      helper,
			KnownHostsFile: "/tmp/known hosts",
		},
		connectTimeout:  2 * time.Second,
		transferTimeout: 2 * time.Second,
	}
}

func TestCmafsHelperProcess(t *testing.T) {
	mode := os.Getenv("CMAFS_HELPER_MODE")
	if mode == "" {
		return
	}
	args := os.Args
	for i, arg := range args {
		if arg == "--" && i+1 < len(args) {
			args = args[i+1:]
			break
		}
	}
	if path := os.Getenv("CMAFS_HELPER_ARGS"); path != "" {
		_ = os.WriteFile(path, []byte(strings.Join(args, "\n")), 0644)
	}

	isRsync := slices.Contains(args, "--protect-args")
	if isRsync {
		runRsyncHelper(mode, args)
		os.Exit(0)
	}

	switch mode {
	case "stat-existing", "changing":
		fmt.Printf("regular file|%d|%d\n", helperDataSize, helperMtime)
	case "stat-missing":
		fmt.Fprintln(os.Stderr, "stat: cannot stat: No such file or directory")
		os.Exit(1)
	case "stat-permission":
		fmt.Fprintln(os.Stderr, "stat: cannot stat: Permission denied")
		os.Exit(1)
	case "stat-host-key":
		fmt.Fprintln(os.Stderr, "Host key verification failed")
		os.Exit(255)
	case "stat-malformed":
		fmt.Println("not-a-stat-record")
	case "stat-timeout":
		time.Sleep(2 * time.Second)
	default:
		fmt.Fprintf(os.Stderr, "unknown helper stat mode %q\n", mode)
		os.Exit(99)
	}
	os.Exit(0)
}

func runRsyncHelper(mode string, args []string) {
	target := args[len(args)-1]
	switch mode {
	case "rsync-success", "changing":
		writeHelperData(target, helperDataSize, helperMtime)
	case "rsync-failure":
		_ = os.WriteFile(target, []byte("partial"), 0644)
		fmt.Fprintln(os.Stderr, "rsync error: simulated failure")
		os.Exit(23)
	case "rsync-timeout", "rsync-cancel":
		time.Sleep(2 * time.Second)
	case "rsync-spawn":
		child := exec.Command("sleep", "30")
		if err := child.Start(); err != nil {
			os.Exit(98)
		}
		_ = os.WriteFile(os.Getenv("CMAFS_HELPER_CHILD_PID"), []byte(strconv.Itoa(child.Process.Pid)), 0644)
		time.Sleep(2 * time.Second)
	default:
		fmt.Fprintf(os.Stderr, "unknown helper rsync mode %q\n", mode)
		os.Exit(99)
	}
}

func writeHelperData(path string, size int, mtime int64) {
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), size), 0644); err != nil {
		os.Exit(97)
	}
	stamp := time.Unix(mtime, 0)
	if err := os.Chtimes(path, stamp, stamp); err != nil {
		os.Exit(97)
	}
}

func TestBackendStatExisting(t *testing.T) {
	backend := helperBackend(t, "stat-existing")
	metadata, err := backend.Stat(context.Background(), "/sstorage/a b'c.grb2")
	require.NoError(t, err)
	assert.True(t, metadata.Regular)
	assert.Equal(t, int64(helperDataSize), metadata.Size)
	assert.Equal(t, int64(helperMtime), metadata.ModTime.Unix())
	args, err := os.ReadFile(os.Getenv("CMAFS_HELPER_ARGS"))
	require.NoError(t, err)
	assert.Contains(t, string(args), "StrictHostKeyChecking=yes")
	assert.Contains(t, string(args), "UserKnownHostsFile=/tmp/known hosts")
	assert.Contains(t, string(args), "'/sstorage/a b'\"'\"'c.grb2'")
}

func TestBackendStatErrorClassification(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		_, err := helperBackend(t, "stat-missing").Stat(context.Background(), "/missing")
		assert.ErrorIs(t, err, fetcher.ErrSourceNotExist)
	})
	for _, mode := range []string{"stat-permission", "stat-host-key", "stat-malformed"} {
		t.Run(mode, func(t *testing.T) {
			_, err := helperBackend(t, mode).Stat(context.Background(), "/data")
			assert.ErrorIs(t, err, fetcher.ErrBackendOperation)
			assert.NotErrorIs(t, err, fetcher.ErrSourceNotExist)
		})
	}
}

func TestBackendStatTimeout(t *testing.T) {
	backend := helperBackend(t, "stat-timeout")
	backend.connectTimeout = 30 * time.Millisecond
	_, err := backend.Stat(context.Background(), "/data")
	assert.ErrorIs(t, err, fetcher.ErrBackendOperation)
}

func TestBackendFetchSuccessAndSafeArguments(t *testing.T) {
	backend := helperBackend(t, "rsync-success")
	target := filepath.Join(t.TempDir(), "target with space.grb2")
	remote := "/sstorage/a b'c.grb2"
	require.NoError(t, backend.Fetch(context.Background(), remote, target))
	info, err := os.Stat(target)
	require.NoError(t, err)
	assert.Equal(t, int64(helperDataSize), info.Size())
	assert.Equal(t, int64(helperMtime), info.ModTime().Unix())
	args, err := os.ReadFile(os.Getenv("CMAFS_HELPER_ARGS"))
	require.NoError(t, err)
	assert.Contains(t, string(args), "--protect-args\n")
	assert.Contains(t, string(args), "server:"+remote+"\n")
	assert.Contains(t, string(args), target)
	lines := strings.Split(string(args), "\n")
	for i, line := range lines {
		if line == "-e" && i+1 < len(lines) {
			assert.NotContains(t, lines[i+1], "server", "rsync -e must not include the destination; rsync appends it")
		}
	}
}

func TestBackendFetchFailureAndTimeout(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		backend := helperBackend(t, "rsync-failure")
		err := backend.Fetch(context.Background(), "/data", filepath.Join(t.TempDir(), "partial"))
		assert.ErrorIs(t, err, fetcher.ErrBackendOperation)
	})
	t.Run("timeout", func(t *testing.T) {
		backend := helperBackend(t, "rsync-timeout")
		backend.transferTimeout = 30 * time.Millisecond
		err := backend.Fetch(context.Background(), "/data", filepath.Join(t.TempDir(), "partial"))
		assert.ErrorIs(t, err, fetcher.ErrBackendOperation)
	})
}

func TestBackendFetchCancellation(t *testing.T) {
	backend := helperBackend(t, "rsync-cancel")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	err := backend.Fetch(ctx, "/data", filepath.Join(t.TempDir(), "partial"))
	assert.ErrorIs(t, err, context.Canceled)
}

func TestBackendFetchCancellationKillsChildProcess(t *testing.T) {
	backend := helperBackend(t, "rsync-spawn")
	backend.transferTimeout = 5 * time.Second
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pidPath := os.Getenv("CMAFS_HELPER_CHILD_PID")
	go func() {
		for range 100 {
			if _, err := os.Stat(pidPath); err == nil {
				cancel()
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		cancel()
	}()
	err := backend.Fetch(ctx, "/data", filepath.Join(t.TempDir(), "partial"))
	require.ErrorIs(t, err, context.Canceled)
	data, err := os.ReadFile(pidPath)
	require.NoError(t, err)
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	require.NoError(t, err)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })

	assert.Eventually(t, func() bool {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return true
		}
		output, psErr := exec.Command("ps", "-o", "stat=", "-p", strconv.Itoa(pid)).Output()
		return psErr != nil || strings.Contains(string(output), "Z")
	}, time.Second, 20*time.Millisecond)
}
