package main

import (
	"encoding/json"
	"os/exec"
)

var cmdShutdown = &Command{
	Name: "guest-shutdown",
	Func: fnShutdown,
}

func init() {
	commands = append(commands, cmdShutdown)
}

func fnShutdown(req *Request) *Response {
	res := &Response{}

	shutdown := struct {
		Mode string `json:"mode"`
	}{}
	ret := struct {
		Id int `json:"-"`
	}{}

	err := json.Unmarshal(req.RawArgs, &shutdown)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	var args []string = []string{"-h"}

	switch shutdown.Mode {
	case "halt":
		args = append(args, "-H")
		break
	case "reboot":
		args = append(args, "-r")
		break
	case "powerdown":
	default:
		args = append(args, "-P")
		break
	}
	args = append(args, "+0", "hypervisor initiated shutdown")
	cmd := exec.Command("shutdown", args...)
	defer cmd.Run()

	return &Response{}
}
