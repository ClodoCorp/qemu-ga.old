package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	flags "github.com/jessevdk/go-flags"
)

type Request struct {
	Execute   string                 `json:"execute"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type Response struct {
	Return interface{} `json:"return"`
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
	var req Request
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

	os.Chdir("/")
	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)

	wait := 10
	for {
		f, err = os.OpenFile(options.Path, os.O_RDWR, os.FileMode(os.ModeCharDevice|0600))
		if err == nil {
			break
		}
		if wait < 0 {
			log.Fatal("Failed to open device:", err)
			os.Exit(1)
		}
		wait -= 1
		time.Sleep(5 * time.Second)
	}

	defer f.Close()
	dec := json.NewDecoder(f)
	dec.UseNumber()
	enc := json.NewEncoder(f)

	for {
		time.Sleep(500 * time.Millisecond)
		err = dec.Decode(&req)
		if err == nil {
			for _, cmd := range commands {
				if cmd.Name == req.Execute && cmd.Func != nil {
					go handle(enc, cmd.Func, req.Arguments)
				}
			}
		}

	}

}

func handle(enc *json.Encoder, fn func(map[string]interface{}) interface{}, args map[string]interface{}) {
	enc.Encode(fn(args))
}
