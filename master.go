package main

import (
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

func master() error {
	var err error

	if err = ioutil.WriteFile("/proc/self/oom_score_adj", []byte("-1000"), 0644); err != nil {
		l.Debug(err.Error())
	}

	if err = ioutil.WriteFile("/proc/self/oom_adj", []byte("-17"), 0644); err != nil {
		l.Debug(err.Error())
	}

	if err = os.Chdir("/"); err != nil {
		l.Debug(err.Error())
	}

	syscall.Umask(0)

	if err = syscall.Setpgid(0, 0); err != nil {
		l.Debug(err.Error())
	}

	if _, err = syscall.Setsid(); err != nil {
		l.Debug(err.Error())
	}

	for _, pid := range getPids("qemu-ga", true) {
		syscall.Kill(pid, syscall.SIGTERM)
	}
	/*
		syscall.Close(0)
		syscall.Close(1)
		syscall.Close(2)
	*/
	ch, err := NewVirtioChannel()
	if err != nil {
		l.Error(err.Error())
	}
	if err = ch.DialTimeout(options.Path, time.Minute); err != nil {
		l.Error(err.Error())
	}
	defer ch.Close()

	ch.Poll()

	return nil
}
