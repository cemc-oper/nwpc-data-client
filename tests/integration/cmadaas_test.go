//go:build integration

package integration

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/tests/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cmadaasTestCase struct {
	name         string
	dataType     string
	startTime    string
	member       string
	expectedPath string
}

var cmadaasTestCases = []cmadaasTestCase{
	{
		name:         "CMA-GFS-GMF",
		dataType:     "cma_gfs_gmf/current/grib2/orig",
		startTime:    "2026111300",
		expectedPath: "/CMADAAS/DATA/NAFP/NMC/GRAPES-GFS-GLB/2026/20261113/Z_NAFP_C_BABJ_20261113000000_P_NWPC-GRAPES-GFS-GLB-02400.grib2",
	},
	{
		name:         "CMA-GEPS",
		dataType:     "cma_geps/current/grib2/orig",
		startTime:    "2026111300",
		member:       "001",
		expectedPath: "/CMADAAS/DATA/NAFP/NMC/GRAPES-GEPS/2026/20261113/Z_NAFP_C_BABJ_20261113000000_P_NWPC-GRAPES-GEPS-GLB-02400-m001.grib2",
	},
	{
		name:         "CMA-TYM",
		dataType:     "cma_tym/current/grib2/orig",
		startTime:    "2026111300",
		expectedPath: "/CMADAAS/DATA/NAFP/NMC/GRAPES-TYM/2026/20261113/Z_NAFP_C_BABJ_20261113000000_P_NWPC-GRAPES-TYM-ACWP-02400.grib2",
	},
	{
		name:         "CMA-REPS 18 UTC next day",
		dataType:     "cma_reps/current/grib2/orig",
		startTime:    "2026111218",
		member:       "001",
		expectedPath: "/CMADAAS/ORIG-DATA/NAFP/CMA-REPS/2026/20261113/Z_NAFP_C_BABJ_20261112180000_P_NWPC-GRAPES-REPS-CN-02400-m001.grib2",
	},
	{
		name:         "CMA-REPS 06 UTC current day",
		dataType:     "cma_reps/current/grib2/orig",
		startTime:    "2026111306",
		member:       "001",
		expectedPath: "/CMADAAS/ORIG-DATA/NAFP/CMA-REPS/2026/20261113/Z_NAFP_C_BABJ_20261113060000_P_NWPC-GRAPES-REPS-CN-02400-m001.grib2",
	},
	{
		name:         "CMA-MESO 3KM",
		dataType:     "cma_meso_3km/current/grib2/orig",
		startTime:    "2026111300",
		expectedPath: "/CMADAAS/DATA/NAFP/GRAPES/RAFS/2026/20261113/Z_NAFP_C_BABJ_20261113000000_P_NWPC-GRAPES-3KM-ORIG-02400.grb2",
	},
	{
		name:         "CMA-MESO 1KM",
		dataType:     "cma_meso_1km/current/grib2/orig",
		startTime:    "2026111300",
		expectedPath: "/CMADAAS/DATA/NAFP/CEMC/MESO_1KM/ORI/2026/20261113/Z_NAFP_C_BABJ_20261113000000_P_CEMC-CMA-MESO-1KM-ORIG-02400.grb2",
	},
}

func requireCMADAASFile(t *testing.T, path string) os.FileInfo {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Skipf("CMADAAS sample not available: %s", path)
	}
	require.NoError(t, err, "stat CMADAAS sample: %s", path)
	require.True(t, info.Mode().IsRegular(), "CMADAAS sample is not a regular file: %s", path)
	require.Positive(t, info.Size(), "CMADAAS sample is empty: %s", path)
	return info
}

func cmadaasArgs(tc cmadaasTestCase) []string {
	args := []string{
		"--data-type=" + tc.dataType,
		"--start-time=" + tc.startTime,
		"--forecast-time=24h",
		"--location-level=archive",
		"--storage-base=/CMADAAS",
	}
	if tc.member != "" {
		args = append(args, "--member="+tc.member)
	}
	return args
}

func runCMADAASFetch(t *testing.T, tc cmadaasTestCase, targetDir string, method string) string {
	t.Helper()
	args := []string{"fetch"}
	args = append(args, cmadaasArgs(tc)...)
	args = append(args,
		"--target-dir="+targetDir,
		"--method="+method,
		"--stable-interval=20ms",
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, testutil.BinaryPath(t), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	require.NotEqual(t, context.DeadlineExceeded, ctx.Err(), "CMADAAS fetch timed out: %s", stderr.String())
	require.NoError(t, err, "CMADAAS fetch failed: %s", stderr.String())
	return strings.TrimSpace(stdout.String())
}

func TestCMADAASMountedArchive(t *testing.T) {
	if _, err := os.Stat("/CMADAAS"); os.IsNotExist(err) {
		t.Skip("CMADAAS mount is not available")
	} else {
		require.NoError(t, err)
	}

	for _, tc := range cmadaasTestCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceInfo := requireCMADAASFile(t, tc.expectedPath)
			args := append([]string{"local"}, cmadaasArgs(tc)...)
			got := testutil.RunCLI(t, args...)
			assert.Equal(t, tc.expectedPath, got)

			targetDir := t.TempDir()
			target := runCMADAASFetch(t, tc, targetDir, "symlink")
			assert.Equal(t, filepath.Join(targetDir, filepath.Base(tc.expectedPath)), target)
			linkInfo, err := os.Lstat(target)
			require.NoError(t, err)
			assert.NotZero(t, linkInfo.Mode()&os.ModeSymlink)
			linkTarget, err := os.Readlink(target)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedPath, linkTarget)
			targetInfo, err := os.Stat(target)
			require.NoError(t, err)
			assert.Equal(t, sourceInfo.Size(), targetInfo.Size())
		})
	}
}

func sha256File(t *testing.T, path string) string {
	t.Helper()
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, file.Close())
	}()
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	require.NoError(t, err)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func TestCMADAASFetchCopy(t *testing.T) {
	tc := cmadaasTestCases[1] // The GEPS sample is the smallest confirmed file.
	sourceInfo := requireCMADAASFile(t, tc.expectedPath)
	targetDir := t.TempDir()
	started := time.Now()
	target := runCMADAASFetch(t, tc, targetDir, "copy")
	elapsed := time.Since(started)

	assert.Equal(t, filepath.Join(targetDir, filepath.Base(tc.expectedPath)), target)
	targetInfo, err := os.Stat(target)
	require.NoError(t, err)
	assert.True(t, targetInfo.Mode().IsRegular())
	assert.Equal(t, sourceInfo.Size(), targetInfo.Size())
	assert.Equal(t, sourceInfo.ModTime().Unix(), targetInfo.ModTime().Unix())
	sourceHash := sha256File(t, tc.expectedPath)
	targetHash := sha256File(t, target)
	assert.Equal(t, sourceHash, targetHash)
	t.Logf("copied %d bytes in %s; sha256=%s", targetInfo.Size(), elapsed, targetHash)
}
