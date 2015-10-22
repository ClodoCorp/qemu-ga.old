package main

import "fmt"

// Command struct contains supported commands
type Command struct {
	Enabled bool                     `json:"enabled"`
	Name    string                   `json:"name"`
	Func    func(*Request) *Response `json:"-"`
	Returns bool                     `json:"success-response"`
}

var commands = []*Command{}

func CmdRun(req *Request) *Response {
	for _, cmd := range commands {
		if cmd.Name == req.Execute && cmd.Func != nil {
			res := cmd.Func(req)
			if cmd.Returns {
				return res
			} else {
				ret := struct{}{}
				return &Response{Return: ret}
			}
		}
	}
	return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("command %s not found", req.Execute)}}
}
