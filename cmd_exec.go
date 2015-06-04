package main

import (
	"encoding/base64"
	"encoding/json"
	"os/exec"
)

var cmdExec = &Command{
	Name: "guest-exec",
	Func: fnExec,
}

func init() {
	commands = append(commands, cmdSync)
}

func fnExec(d map[string]interface{}) interface{} {
	var result struct {
		ExitCode int
		Output   string
	}
	str, _ := (d["command"].(json.Number)).Int64()
	cmdline, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return &Response{
			Return: result{
				ExitCode: 1,
				Output:   err.Error(),
			},
		}
	}
	output, err := exec.Command("sh", "-c", "'", string(data), "'").CombinedOutput()
	var ret result

	if err != nil {
		ret.ExitCode = 1
	}
	ret.Output = base64.StdEncoding.EncodeToString(output)

	return &Response{
		Return: ret,
	}
}
