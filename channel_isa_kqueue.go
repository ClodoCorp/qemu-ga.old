// +build freebsd openbsd netbsd darwin

package main

import (
	"encoding/json"
	"fmt"
	"syscall"
)

func (ch *IsaChannel) Poll() error {
	var err error

	ch.fd = int(ch.f.Fd())

	if err = syscall.SetNonblock(ch.fd, true); err != nil {
		return err
	}

	ch.pfd, err = syscall.Kqueue()
	if err != nil {
		return err
	}

	ctlEvent := syscall.Kevent_t{
		Ident:  uint64(ch.fd),
		Filter: syscall.EVFILT_VNODE | syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}
	timeout := syscall.Timespec{
		Sec:  0,
		Nsec: 0,
	}
	events := make([]syscall.Kevent_t, 32)
	if _, err = syscall.Kevent(ch.pfd, events, nil, nil); err != nil {
		return err
	}

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

	buffer := make([]byte, 4*1024)
	var n int
	var req Request
	for {

		nevents, err := syscall.Kevent(ch.pfd, nil, events, &timeout)
		switch err {
		case nil:
			for ev := 0; ev < nevents; ev++ {
				n, err = syscall.Read(int(events[ev].Ident), buffer)
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
			return err
		}
	}

	//	return fmt.Errorf("isa channel poll failed")
}

func (ch *IsaChannel) Close() error {
	if err := syscall.Close(ch.pfd); err != nil {
		return err
	}
	return ch.f.Close()
}
