package main

import (
	"encoding/base64"
	"fmt"
	"os/exec"
)

var cmdPasswd = &Command{
	Name: "guest-set-user-password",
	Func: fnPasswd,
}

func init() {
	commands = append(commands, cmdPasswd)
}

func fnPasswd(d map[string]interface{}) interface{} {
	type request struct {
		User    string `json:"username"`
		Hash    string `json:"password"`
		Crypted bool   `json:"crypted"`
	}

	user := d["username"].(string)
	crypted := d["crypted"].(bool)
	hash, err := base64.StdEncoding.DecodeString(d["password"].(string))
	if err != nil {
		return &Response{}
	}

	args := []string{}

	if crypted {
		args = append(args, "-e")
	}

	cmd := exec.Command("chpasswd", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return &Response{}
	}

	err = cmd.Start()
	if err != nil {
		return &Response{}
	}

	arg := fmt.Sprintf("%s:%s", user, hash)
	_, err = stdin.Write([]byte(arg))
	if err != nil {
		return &Response{}
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return &Response{}
	}

	return &Response{}
}
