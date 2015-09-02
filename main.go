package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

func main() {
	var req Request
	var f io.ReadWriteCloser

	parser := flags.NewParser(&options, flags.PrintErrors)
	_, err := parser.Parse()
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}

	if options.Help {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	wait := 5
	for {
		f, err = os.OpenFile(options.Path, os.O_RDWR|os.O_APPEND, os.FileMode(os.ModeCharDevice|0600))
		if err == nil {
			break
		}
		if wait < 0 {
			log.Fatal("Failed to open device:", err)
		}
		wait -= 1
		time.Sleep(5 * time.Second)
	}

	dec := json.NewDecoder(f)
	dec.UseNumber()
	enc := json.NewEncoder(f)

	for {
		time.Sleep(500 * time.Millisecond)
		dec.Decode(&req)
		for _, cmd := range commands {
			if cmd.Name == req.Execute && cmd.Func != nil {
				enc.Encode(cmd.Func(req.Arguments))
			}
		}
		if finish {
			break
		}
	}

	f.Close()
	os.Exit(0)

}
