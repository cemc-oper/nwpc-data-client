package finder

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCandidatePathsCmafsUsesPOSIXPath(t *testing.T) {
	config := dataconfig.DataConfig{
		FileName: "data.grb2",
		Paths: []dataconfig.PathItem{{
			PathType: "cmafs", LevelType: "storage", Path: "/sstorage/COMMONDATA",
		}},
	}
	candidates := CandidatePaths(config, nil)
	require.Len(t, candidates, 1)
	assert.Equal(t, path.Join("/sstorage/COMMONDATA", "data.grb2"), candidates[0].Path)
	assert.Equal(t, "storage", candidates[0].LevelType)
}

func TestFindLocalFileSkipsCmafs(t *testing.T) {
	dir := t.TempDir()
	remoteLookingPath := filepath.Join(dir, "data.grb2")
	require.NoError(t, os.WriteFile(remoteLookingPath, []byte("data"), 0644))
	config := dataconfig.DataConfig{
		Default:  "NOTFOUND",
		FileName: "data.grb2",
		Paths: []dataconfig.PathItem{{
			PathType: "cmafs", LevelType: "storage", Path: dir,
		}},
	}

	result := FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)
	assert.Equal(t, "NOTFOUND", result.Path)
}
