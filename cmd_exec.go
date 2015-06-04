package main

import (
	"encoding/base64"
	"os/exec"
)

var cmdExec = &Command{
	Name: "guest-exec",
	Func: fnExec,
}

func init() {
	commands = append(commands, cmdExec)
}

func fnExec(d map[string]interface{}) interface{} {
	type result struct {
		ExitCode int
		Output   string
	}
	var ret result

	str, _ := (d["command"]).(string)
	cmdline, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		ret.ExitCode = 1
		ret.Output = base64.StdEncoding.EncodeToString([]byte(err.Error()))
		return &Response{
			Return: ret,
		}
	}

	output, err := exec.Command("sh", "-c", "'"+string(cmdline)+"'").CombinedOutput()

	if err != nil {
		ret.ExitCode = 1
	}
	ret.Output = base64.StdEncoding.EncodeToString(output)

	return &Response{
		Return: ret,
	}
}
