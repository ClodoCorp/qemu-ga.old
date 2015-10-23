package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
)

var cmdPasswd = &Command{
	Name:    "guest-set-user-password",
	Func:    fnPasswd,
	Enabled: true,
}

func init() {
	commands = append(commands, cmdPasswd)
}

func fnPasswd(req *Request) *Response {
	res := &Response{}
	pwd := struct {
		User    string `json:"username"`
		Passwd  string `json:"password"`
		Crypted bool   `json:"crypted"`
	}{}

	ret := struct {
		id int `json:"-"`
	}{}

	err := json.Unmarshal(req.RawArgs, &pwd)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	passwd, err := base64.StdEncoding.DecodeString(pwd.Passwd)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{}

	if pwd.Crypted {
		args = append(args, "-e")
	}

	cmd := exec.Command("chpasswd", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	err = cmd.Start()
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	arg := fmt.Sprintf("%s:%s", pwd.User, passwd)
	_, err = stdin.Write([]byte(arg))
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	res.Return = ret
	return res
}
