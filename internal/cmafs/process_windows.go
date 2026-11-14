//go:build windows

package cmafs

import "os/exec"

// The cmafs backend is operational on HPC Unix systems. Windows retains the
// default CommandContext process termination so cross-platform builds work.
func configureProcessGroup(cmd *exec.Cmd) {}
