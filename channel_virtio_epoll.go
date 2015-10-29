// +build linux

package main

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"
)

func (ch *VirtioChannel) Poll() error {
	var err error

	ch.fd = int(ch.f.Fd())

	if err = syscall.SetNonblock(ch.fd, true); err != nil {
		return err
	}

	ch.pfd, err = syscall.EpollCreate1(0)
	if err != nil {
		return err
	}

	ctlEvent := syscall.EpollEvent{Events: syscall.EPOLLIN | syscall.EPOLLHUP, Fd: int32(ch.fd)}
	if err = syscall.EpollCtl(ch.pfd, syscall.EPOLL_CTL_ADD, ch.fd, &ctlEvent); err != nil {
		return err
	}
	events := make([]syscall.EpollEvent, 32)

	chErr := make(chan error)

	go func() {
		buffer := make([]byte, 4*1024)
		var n int
		var req Request
		for {
			nevents, err := syscall.EpollWait(ch.pfd, events, -1)
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
				chErr <- err
				return
			}
		}
	}()

	go func() {
		var n int
		timer := time.NewTimer(time.Minute * 5)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				chErr <- fmt.Errorf("timeout waiting for command")
				return
			case req := <-ch.req:
				timer.Reset(time.Minute * 5)
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

	for {
		select {
		case err := <-chErr:
			return err
		}
	}
}
