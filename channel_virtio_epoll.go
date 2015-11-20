// +build linux

package main

import (
	"encoding/json"
	"fmt"

	"github.com/vtolstov/qemu-ga/qga"

	"golang.org/x/sys/unix"
)

func (ch *VirtioChannel) Poll() error {
	var err error

	ch.fd = int(ch.f.Fd())

	if err = unix.SetNonblock(ch.fd, true); err != nil {
		return err
	}

	ch.pfd, err = unix.EpollCreate(1)
	if err != nil {
		return err
	}

	ctlEvent := unix.EpollEvent{Events: unix.EPOLLIN | unix.EPOLLHUP, Fd: int32(ch.fd)}
	if err = unix.EpollCtl(ch.pfd, unix.EPOLL_CTL_ADD, ch.fd, &ctlEvent); err != nil {
		return err
	}
	events := make([]unix.EpollEvent, 32)

	chErr := make(chan error, 1)
	defer close(chErr)
	done := make(chan struct{}, 1)

	go func() {

		buffer := make([]byte, 4*1024)
		var n int
		var req qga.Request
		for {
			nevents, err := unix.EpollWait(ch.pfd, events, 1000*60*5)
			switch err {
			case nil:
				if nevents == 0 {
					done <- struct{}{}
					chErr <- fmt.Errorf("timeout waiting for command")
					return
				}
				for ev := 0; ev < nevents; ev++ {
					n, err = unix.Read(int(events[ev].Fd), buffer)
					if err == nil {
						err = json.Unmarshal(buffer[:n], &req)
						if err == nil {
							ch.req <- &req
						} else {
							ch.res <- &qga.Response{Error: &qga.Error{Code: -1, Desc: fmt.Sprintf("invalid request %s", err.Error())}}
						}
					}
				}
			case unix.EINTR:
				continue
			default:
				chErr <- err
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case req := <-ch.req:
				ch.res <- qga.CmdRun(req)
			case res := <-ch.res:
				buffer, err := json.Marshal(res)
				buffer = append(buffer, []byte("\n")...)
				if err == nil {
					_, err = unix.Write(ch.fd, buffer)
					fmt.Printf(err.Error())
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
