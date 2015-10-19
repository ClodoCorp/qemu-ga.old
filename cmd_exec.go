package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os/exec"
)

var cmdExec = &Command{
	Name: "guest-exec",
	Func: fnExec,
}

func init() {
	commands = append(commands, cmdExec)
}

func fnExec(m json.RawMessage) json.RawMessage {
	res := struct {
		Return struct {
			ExitCode int
			Output   string
		} `json:"return"`
	}{}

	req := struct {
		Command string `json:"command"`
	}{}

	err := json.Unmarshal(m, &req)
	if err != nil {
		log.Printf("RRR %s\n", err.Error())
	}

	cmdline, err := base64.StdEncoding.DecodeString(req.Command)
	if err != nil {
		res.Return.ExitCode = 1
		res.Return.Output = base64.StdEncoding.EncodeToString([]byte(err.Error()))
		buf, err := json.Marshal(res)
		if err != nil {
			log.Printf("RRR %s\n", err.Error())
		}
		return json.RawMessage(buf)
	}

	output, err := exec.Command("sh", "-c", string(cmdline)).CombinedOutput()

	if err != nil {
		res.Return.ExitCode = 1
	}
	res.Return.Output = base64.StdEncoding.EncodeToString(output)

	buf, err := json.Marshal(res)
	if err != nil {
		log.Printf("RRR %s\n", err.Error())
	}
	return json.RawMessage(buf)
}
