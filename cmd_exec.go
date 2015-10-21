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
	commands = append(commands, cmdExec)
}

func fnExec(req *Request) *Response {
	res := &Response{}

	ex := struct {
		ExitCode int
		Output   string
	}{}

	cmd := struct {
		Command string `json:"command"`
	}{}

	err := json.Unmarshal(req.RawArgs, &cmd)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	if cmd.Command == "" {
		res.Error = &Error{Code: -1, Desc: "empty command to guest-exec"}
		return res
	}
	cmdline, err := base64.StdEncoding.DecodeString(cmd.Command)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	output, err := exec.Command("sh", "-c", string(cmdline)).CombinedOutput()
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	ex.Output = base64.StdEncoding.EncodeToString(output)
	ex.ExitCode = 0
	res.Return = ex
	return res
}
