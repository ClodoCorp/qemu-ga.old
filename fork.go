// +build !windows

package main

import (
	"os/exec"
	"syscall"
)

func fork() error {
	cmd := exec.Command("qemu-ga")
	cmd.Dir = "/"
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.ExtraFiles = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: false, Setpgid: true}

	return cmd.Start()
}
