// +build linux

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"time"
)

type VirtioChannel struct {
	f    *os.File
	fd   int
	epfd int
	req  chan *Request
	res  chan *Response
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
		if f, err = os.OpenFile(path, os.O_RDWR|syscall.O_NONBLOCK|syscall.O_ASYNC|syscall.O_CLOEXEC|syscall.O_NDELAY, os.FileMode(os.ModeCharDevice|0600)); err == nil {
			ch.f = f
			ch.req = make(chan *Request)
			ch.res = make(chan *Response, 1)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("virtio channel failed to connect")
}

func (ch *VirtioChannel) Poll() error {
	var err error

	ch.fd = int(ch.f.Fd())

	if err = syscall.SetNonblock(ch.fd, true); err != nil {
		return err
	}

	ch.epfd, err = syscall.EpollCreate1(0)
	if err != nil {
		return err
	}

	ctlEvent := syscall.EpollEvent{Events: syscall.EPOLLIN | syscall.EPOLLHUP, Fd: int32(ch.fd)}
	if err = syscall.EpollCtl(ch.epfd, syscall.EPOLL_CTL_ADD, ch.fd, &ctlEvent); err != nil {
		return err
	}
	events := make([]syscall.EpollEvent, 32)

	go func() {
		var n int
		for {
			select {
			case req := <-ch.req:
				ch.res <- CmdRun(req)
			case res := <-ch.res:
				buffer, err := json.Marshal(res)
				buffer = append(buffer, []byte("\n")...)
				if err == nil {
					n, err = syscall.Write(ch.fd, buffer)
					_ = n
					_ = err
				} else {
					fmt.Printf(err.Error())
				}
			}
		}
	}()

	go func() {
		buffer := make([]byte, 4*1024)
		var n int
		var req Request
		for {
			nevents, err := syscall.EpollWait(ch.epfd, events, -1)
			switch err {
			case nil:
				for ev := 0; ev < nevents; ev++ {
					n, err = syscall.Read(int(events[ev].Fd), buffer)
					if err == nil {
						err = json.Unmarshal(buffer[:n], &req)
						if err == nil {
							ch.req <- &req
						}
					}
				}
			case syscall.EINTR:
				continue
			default:
				break
			}
		}
	}()
	select {}
	return fmt.Errorf("channel virtio poll failed")
}

func (ch *VirtioChannel) Close() error {
	if err := syscall.Close(ch.epfd); err != nil {
		return err
	}
	return ch.f.Close()
}
