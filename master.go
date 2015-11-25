package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func master() error {
	var err error

	master := true
	if os.Getenv("master") == "false" {
		master = false
	}

	if err = os.Chdir("/"); err != nil {
		l.Debug(err.Error())
	}

	unix.Umask(0)

	if master {
		if err = ioutil.WriteFile("/proc/self/oom_score_adj", []byte("-1000"), 0644); err != nil {
			l.Debug(err.Error())
		}

		if err = ioutil.WriteFile("/proc/self/oom_adj", []byte("-17"), 0644); err != nil {
			l.Debug(err.Error())
		}

		if err = unix.Setpgid(0, 0); err != nil {
			l.Debug(err.Error())
		}

		if _, err = unix.Setsid(); err != nil {
			l.Debug(err.Error())
		}

		for _, pid := range getPids("qemu-ga", true) {
			unix.Kill(pid, unix.SIGTERM)
		}
	}

	/*
		syscall.Close(0)
		syscall.Close(1)
		syscall.Close(2)
	*/

	if master {
		stdOut := bytes.NewBuffer(nil)
		stdErr := bytes.NewBuffer(nil)
		for {
			cmd := exec.Command("qemu-ga")
			cmd.Dir = "/"
			cmd.Env = append(cmd.Env, "master", "false")
			cmd.Stdin = nil
			cmd.Stdout = stdOut
			cmd.Stderr = stdErr
			cmd.ExtraFiles = nil
			cmd.Run()
			stdErr.Reset()
			stdOut.Reset()
		}
	} else {
		for {
			if err = slave(); err != nil {
				l.Error(err.Error())
			}
		}
	}

}
