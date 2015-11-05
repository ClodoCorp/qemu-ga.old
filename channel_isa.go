package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/vtolstov/qemu-ga/qga"
)

type IsaChannel struct {
	f   *os.File
	fd  int
	pfd int
	req chan *qga.Request
	res chan *qga.Response
}

func NewIsaChannel() (*IsaChannel, error) {
	return &IsaChannel{}, nil
}

func (ch *IsaChannel) DialTimeout(path string, timeout time.Duration) error {
	var f *os.File
	var err error

	select {
	case <-time.After(timeout):
		return fmt.Errorf("isa channel dial timeout: %s", path)
	default:
		if f, err = os.OpenFile(path, os.O_RDWR|syscall.O_NONBLOCK|syscall.O_ASYNC|syscall.O_CLOEXEC|syscall.O_NDELAY, os.FileMode(os.ModeCharDevice|0600)); err == nil {
			ch.f = f
			ch.req = make(chan *qga.Request)
			ch.res = make(chan *qga.Response, 1)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("isa channel failed to connect")
}
