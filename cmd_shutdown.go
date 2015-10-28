package main

import (
	"encoding/json"
	"os/exec"
)

var cmdShutdown = &Command{
	Name:    "guest-shutdown",
	Func:    fnShutdown,
	Enabled: true,
}

func init() {
	commands = append(commands, cmdShutdown)
}

func fnShutdown(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		Mode string `json:"mode"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{"-h"}

	switch reqData.Mode {
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

	return res
}
