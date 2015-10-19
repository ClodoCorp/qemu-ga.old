package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var mt sync.Mutex

type Request struct {
	Execute   string          `json:"execute"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

func (r *Request) Write(b []byte) (int, error) {
	fmt.Printf("%s\n", b)
	return len(b), nil
}

var finish bool = false

var ppid int

func init() {
	ppid = os.Getpid()
}

func main() {
	var event syscall.EpollEvent
	var events [32]syscall.EpollEvent

	//var req Request
	var f io.ReadWriteCloser

	parser := flags.NewParser(&options, flags.PrintErrors)
	_, err := parser.Parse()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	if options.Help {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	if options.Version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	/*
		parent := os.Getenv("PARENT")
		if parent != "" {
			if pid, err := strconv.Atoi(parent); err == nil {
				if proc, err := os.FindProcess(pid); err == nil {
					err = proc.Kill()
					if err != nil {
						log.Printf(err.Error())
						os.Exit(1)
					}
				}
			}
		} else {
			pids := getPids("qemu-ga", true)
			err = syscall.Setpgid(0, 0)
			if err != nil {
				log.Printf("setpgid err: %s\n", err.Error())
			} else {
				syscall.Setsid()
				for _, pid := range pids {
					syscall.Kill(pid, syscall.SIGTERM)
				}
				time.Sleep(10 * time.Second)
			}
		}
	*/
	os.Chdir("/")
	//	syscall.Close(0)
	//	syscall.Close(1)
	//	syscall.Close(2)

	wait := 10
	for {
		f, err = os.OpenFile(options.Path, os.O_RDWR|syscall.O_NONBLOCK|syscall.O_ASYNC|syscall.O_CLOEXEC, os.FileMode(os.ModeCharDevice|0600))
		if err == nil {
			break
		}
		if wait < 0 {
			log.Fatal("Failed to open device: ", err)
			os.Exit(1)
		}
		wait -= 1
		time.Sleep(5 * time.Second)
	}

	defer f.Close()
	fd := int(f.(*os.File).Fd())
	/*
		if err = syscall.SetNonblock(fd, true); err != nil {
			log.Fatal("Setnonblock: ", err)
			os.Exit(1)
		}
	*/
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatalf("EpollCreate: ", err)
		os.Exit(1)
	}
	defer syscall.Close(epfd)

	event.Events = syscall.EPOLLIN | syscall.EPOLLHUP
	event.Fd = int32(fd)
	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
		log.Fatalf("EpollCtl: ", err)
		os.Exit(1)
	}

	for {
		nevents, err := syscall.EpollWait(epfd, events[:], -1)
		if err != nil && err != syscall.EINTR {
			panic(err)

		}
		log.Printf("epoolwait\n")
		for ev := 0; ev < nevents; ev++ {
			go read(int(events[ev].Fd))
		}

	}

}

func read(fd int) {
	var buf [4096]byte
	var req Request

	n, err := syscall.Read(fd, buf[:])

	if err != nil && err != syscall.EAGAIN && err != io.EOF {
		log.Printf("zzz %s\n", err.Error())
	}
	if n < 1 {
		return
	}
	err = json.Unmarshal(buf[:n], &req)
	if err == nil {
		for _, cmd := range commands {
			if req.Execute == "guest-sync" {
				log.Printf("%s\n", req.Execute)
			}
			if cmd.Name == req.Execute && cmd.Func != nil {
				go write(fd, cmd.Func, req.Arguments)
			}
		}
	} else {
		log.Printf("rr %s\n", err.Error())
	}

}

func write(fd int, fn func(json.RawMessage) json.RawMessage, args json.RawMessage) {
	m := fn(args)

	buf, err := m.MarshalJSON()
	if err != nil {
		return
	}
	if len(buf) > 0 {
		buf = append(buf, []byte("\n")...)
	}
	mt.Lock()
	_, err = syscall.Write(fd, buf)
	mt.Unlock()
	return
}
