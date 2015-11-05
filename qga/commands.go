package qga

import (
	"fmt"
)

// Command struct contains supported commands
type Command struct {
	Enabled bool                     `json:"enabled"`
	Name    string                   `json:"name"`
	Func    func(*Request) *Response `json:"-"`
	Returns bool                     `json:"success-response"`
}

var commands = []*Command{}

func RegisterCommand(cmd *Command) {
	commands = append(commands, cmd)
}

func ListCommands() []*Command {
	return commands
}

func CmdRun(req *Request) *Response {
	if req == nil || req.Execute == "" {
		return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("invalid command")}}
	}
	for _, cmd := range commands {
		if cmd.Name == req.Execute && cmd.Func != nil {
			res := cmd.Func(req)
			if cmd.Returns || res.Error != nil {
				return res
			} else {
				return &Response{Return: struct{}{}}
			}
		}
	}
	return &Response{Error: &Error{Class: "CommandNotFound", Desc: fmt.Sprintf("command %s not found", req.Execute)}}
}
