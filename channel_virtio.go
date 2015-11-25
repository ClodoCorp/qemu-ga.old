// +build !windows

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/vtolstov/qemu-ga/qga"

	"golang.org/x/sys/unix"
)

type VirtioChannel struct {
	f   *os.File
	fd  int
	pfd int
	req chan *qga.Request
	res chan *qga.Response
}

func NewVirtioChannel() (*VirtioChannel, error) {
	return &VirtioChannel{}, nil
}

func (ch *VirtioChannel) DialTimeout(path string, timeout time.Duration) error {
	var f *os.File
	var err error

	select {
	case <-time.After(timeout):
		return fmt.Errorf("virtio channel dial timeout: %s", path)
	default:
		if f, err = os.OpenFile(path, os.O_RDWR|unix.O_NONBLOCK|unix.O_ASYNC|unix.O_CLOEXEC|unix.O_NDELAY, os.FileMode(os.ModeCharDevice|0600)); err == nil {
			ch.f = f
			ch.req = make(chan *qga.Request)
			ch.res = make(chan *qga.Response, 1)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("virtio channel failed to connect")
}

func (ch *VirtioChannel) Close() error {
	if err := unix.Close(ch.pfd); err != nil {
		return err
	}
	close(ch.req)
	close(ch.res)
	return ch.f.Close()
}
