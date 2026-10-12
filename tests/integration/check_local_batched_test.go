//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/tests/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCheckLocalBatchedSmoke runs one batched-mode check against real
// operational CMA-GFS data. It skips when the data is not available (i.e.
// anywhere outside the CMA HPC filesystems).
func TestCheckLocalBatchedSmoke(t *testing.T) {
	checkDateTime := dateHour(yesterday(1), "00")
	dataType := "cma_gfs_gmf/current/grib2/orig"
	forecastTimes := []string{"3h", "6h"}

	// Both files must exist, otherwise the batched loop would poll forever.
	var expectedPaths []string
	for _, ft := range []string{"003", "006"} {
		expectedPaths = append(expectedPaths, fmt.Sprintf(
			"/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/%s/data/output/grib2_orig/gmf.gra.%s%s.grb2",
			checkDateTime, checkDateTime, ft))
	}
	for _, p := range expectedPaths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Skipf("data not available: %s", p)
		}
	}

	binPath := testutil.BinaryPath(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binPath, "check", "local", checkDateTime,
		"--data-type="+dataType,
		"--location-level=runtime",
		"--check-mode=batched",
		"--check-interval=2s",
		"--check-interval-jitter=0",
		"--max-check-count=10",
		"--execute-command=echo ready: {{.FilePath}}",
	)
	cmd.Stdin = strings.NewReader(strings.Join(forecastTimes, " "))

	output, err := cmd.CombinedOutput()
	require.NotEqual(t, context.DeadlineExceeded, ctx.Err(), "batched check timed out: %s", string(output))
	require.NoError(t, err, "batched check local failed: %s", string(output))
	for _, p := range expectedPaths {
		assert.Contains(t, string(output), "ready: "+p)
	}
}
