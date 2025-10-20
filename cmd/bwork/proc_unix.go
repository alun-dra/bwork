//go:build !windows

package main

import (
	"os/exec"
	"syscall"
	"time"
)

func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessTree(pid int) {
	_ = syscall.Kill(-pid, syscall.SIGTERM)
	time.Sleep(300 * time.Millisecond)
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}
