package finder

import (
	"path/filepath"
	"testing"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/stretchr/testify/assert"
)

func TestCandidatePathsPreserveOrderAndMetadata(t *testing.T) {
	config := dataconfig.DataConfig{
		FileName:  "single.grb2",
		FileNames: []string{"first.grb2", "second.grb2"},
		Paths: []dataconfig.PathItem{
			{PathType: "local", LevelType: "runtime", Path: "/runtime"},
			{PathType: "local", LevelType: "archive", Path: "/archive"},
		},
	}

	candidates := CandidatePaths(config, []string{"archive"})

	if assert.Len(t, candidates, 3) {
		assert.Equal(t, filepath.Join("/archive", "first.grb2"), candidates[0].Path)
		assert.Equal(t, filepath.Join("/archive", "second.grb2"), candidates[1].Path)
		assert.Equal(t, filepath.Join("/archive", "single.grb2"), candidates[2].Path)
	}
	for _, candidate := range candidates {
		assert.Equal(t, "local", candidate.PathType)
		assert.Equal(t, "archive", candidate.LevelType)
	}
}
