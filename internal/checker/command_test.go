package checker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/tests/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCommandTemplatesSingle(t *testing.T) {
	templates, err := buildCommandTemplates("echo {{.FilePath}}", nil)
	require.NoError(t, err)
	assert.Len(t, templates, 1)
}

func TestBuildCommandTemplatesList(t *testing.T) {
	templates, err := buildCommandTemplates("", []string{"echo {{.FilePath}}", "cat {{.FilePath}}"})
	require.NoError(t, err)
	assert.Len(t, templates, 2)
}

func TestBuildCommandTemplatesMutualExclusion(t *testing.T) {
	_, err := buildCommandTemplates("echo {{.FilePath}}", []string{"cat {{.FilePath}}"})
	assert.Error(t, err)
}

func TestBuildCommandTemplatesEmpty(t *testing.T) {
	templates, err := buildCommandTemplates("", nil)
	require.NoError(t, err)
	assert.Empty(t, templates)
}

func TestBuildCommandTemplatesFilePathTemplate(t *testing.T) {
	templates, err := buildCommandTemplates("echo {{.FilePath}}", nil)

	assert.NoError(t, err)
	assert.Len(t, templates, 1)

	var builder strings.Builder
	err = templates[0].Execute(&builder, CheckerTemplateVariable{FilePath: "/some/data.grb2"})
	assert.NoError(t, err)
	assert.Equal(t, "echo /some/data.grb2", builder.String())
}

func TestBuildCommandTemplatesInvalidTemplate(t *testing.T) {
	_, err := buildCommandTemplates("echo {{.FilePath", nil)

	assert.Error(t, err)
}

// TestCommandHelperProcess is executed in a child copy of the test binary on
// non-Windows platforms. The explicit -- argument keeps it inactive normally.
func TestCommandHelperProcess(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--" || i+1 >= len(os.Args) {
			continue
		}
		switch os.Args[i+1] {
		case "wait":
			if marker := os.Getenv("NWPC_CHECKER_TEST_COMMAND_STARTED"); marker != "" {
				if err := os.WriteFile(marker, []byte("started"), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "write command helper marker: %v\n", err)
					os.Exit(8)
				}
			}
			time.Sleep(2 * time.Second)
		case "brief":
			time.Sleep(100 * time.Millisecond)
		case "fail":
			os.Exit(7)
		}
		return
	}
}

func commandHelperCommand(t *testing.T, mode string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		// Killing cmd.exe does not kill a child copy of the test binary on
		// Windows. Such an orphan keeps checker.test.exe locked and can make
		// the go tool fail with Access denied while cleaning its build cache.
		// Use PowerShell for waits so no test-binary image remains running.
		switch mode {
		case "success":
			return "exit 0"
		case "fail":
			return "exit 7"
		case "wait":
			script := `
if ($env:NWPC_CHECKER_TEST_COMMAND_STARTED) {
    [System.IO.File]::WriteAllText($env:NWPC_CHECKER_TEST_COMMAND_STARTED, 'started')
}
Start-Sleep -Seconds 2
`
			return testutil.PowerShellEncodedCommand(script)
		case "brief":
			return testutil.PowerShellEncodedCommand("Start-Sleep -Milliseconds 100")
		default:
			t.Fatalf("unknown Windows command helper mode: %s", mode)
		}
	}
	// cmd.exe /C cannot reliably start a quoted executable through the command
	// library's argument escaping. Resolve the helper by name instead; this also
	// keeps the command working when the test binary directory contains spaces.
	t.Setenv("PATH", filepath.Dir(os.Args[0])+string(os.PathListSeparator)+os.Getenv("PATH"))
	parts := []string{
		filepath.Base(os.Args[0]),
		"-test.run=TestCommandHelperProcess",
		"--",
		mode,
	}
	return strings.Join(parts, " ")
}

func testCommandTemplate(t *testing.T, command string) *template.Template {
	t.Helper()
	templates, err := buildCommandTemplates(command, nil)
	require.NoError(t, err)
	return templates[0]
}

func runTestCommand(t *testing.T, ctx context.Context, timeout time.Duration, command string) error {
	t.Helper()
	return runCommand(ctx, timeout, testCommandTemplate(t, command), time.Now(), 0, "", "data.grb2")
}

func TestRunCommandDefaultTimeoutSuccess(t *testing.T) {
	err := runTestCommand(t, context.Background(), 30*time.Minute, commandHelperCommand(t, "success"))
	assert.NoError(t, err)
}

func TestRunCommandNonZeroExit(t *testing.T) {
	err := runTestCommand(t, context.Background(), 30*time.Minute, commandHelperCommand(t, "fail"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exit code is not 0: 7")
}

func TestRunCommandShortTimeout(t *testing.T) {
	err := runTestCommand(t, context.Background(), 50*time.Millisecond, commandHelperCommand(t, "wait"))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCommandTimeout)
	assert.NotErrorIs(t, err, context.Canceled)
	assert.Contains(t, err.Error(), "command timed out after 50ms")
}

func TestRunCommandParentCancellation(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "started")
	t.Setenv("NWPC_CHECKER_TEST_COMMAND_STARTED", marker)
	ctx, cancel := context.WithCancel(context.Background())
	commandTemplate := testCommandTemplate(t, commandHelperCommand(t, "wait"))
	done := make(chan error, 1)
	go func() {
		done <- runCommand(ctx, 30*time.Minute, commandTemplate, time.Now(), 0, "", "data.grb2")
	}()
	require.Eventually(t, func() bool {
		_, err := os.Stat(marker)
		return err == nil
	}, time.Second, 5*time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
		assert.NotErrorIs(t, err, ErrCommandTimeout)
	case <-time.After(time.Second):
		t.Fatal("command did not return after parent cancellation")
	}
}

func TestRunCommandParentDeadlineIsNotCommandTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := runTestCommand(t, ctx, 30*time.Minute, commandHelperCommand(t, "wait"))
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.NotErrorIs(t, err, ErrCommandTimeout)
}

func TestRunCommandZeroTimeoutDisablesDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := runTestCommand(t, ctx, 0, commandHelperCommand(t, "brief"))
	assert.NoError(t, err)
}
