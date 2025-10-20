//go:build windows

package main

import (
	"os/exec"
	"strconv"
)

func setProcessGroup(cmd *exec.Cmd) {
	// no-op en Windows
}

func killProcessTree(pid int) {
	_ = exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
}
