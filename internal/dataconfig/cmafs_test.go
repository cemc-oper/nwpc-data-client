package dataconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLintConfigContentCmafs(t *testing.T) {
	valid := `_config_version: 3
file_name: a.grb2
paths:
  - type: local
    level: runtime
    path: /some/runtime
  - type: cmafs
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC
`
	require.NoError(t, LintConfigContent(valid))

	relative := `_config_version: 3
file_name: a.grb2
paths:
  - type: cmafs
    level: storage
    path: relative/path
`
	err := LintConfigContent(relative)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cmafs path must be absolute")

	newline := "_config_version: 3\nfile_name: a.grb2\npaths:\n  - type: cmafs\n    level: storage\n    path: \"/some/path\\nnext\"\n"
	err = LintConfigContent(newline)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not contain newlines")
}
