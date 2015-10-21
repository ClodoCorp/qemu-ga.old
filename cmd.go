package main

import "fmt"

// Command struct contains supported commands
type Command struct {
	Name string
	Func func(*Request) *Response
}

var commands = []*Command{}

func CmdRun(req *Request) *Response {
	for _, cmd := range commands {
		if cmd.Name == req.Execute && cmd.Func != nil {
			return cmd.Func(req)
		}
	}
	return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("command %s not found", req.Execute)}}
}
