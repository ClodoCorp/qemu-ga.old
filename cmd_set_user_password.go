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
	res := &Response{Id: req.Id}

	reqData := struct {
		User    string `json:"username"`
		Passwd  string `json:"password"`
		Crypted bool   `json:"crypted"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	passwd, err := base64.StdEncoding.DecodeString(reqData.Passwd)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{}

	if reqData.Crypted {
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

	arg := fmt.Sprintf("%s:%s", reqData.User, passwd)
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

	return res
}
