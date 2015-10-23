package main

import (
	"encoding/json"
	"os"
)

var cmdFileOpen = &Command{
	Name:    "guest-file-open",
	Func:    fnFileOpen,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdFileOpen)
}

func fnFileOpen(req *Request) *Response {
	res := &Response{Id: req.Id}

	file := struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
	}{}

	err := json.Unmarshal(req.RawArgs, &file)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	var flag int
	for _, s := range file.Mode {
		switch s {
		case 'a':
			flag = flag | os.O_APPEND | os.O_CREATE | os.O_WRONLY
		case '+':
			flag = flag | os.O_RDWR
		case 'w':
			flag = flag | os.O_TRUNC | os.O_WRONLY
		case 'r':
			flag = flag | os.O_RDONLY
		}
	}

	if f, err := os.OpenFile(file.Path, flag, os.FileMode(0600)); err == nil {
		fd := int(f.Fd())
		openFiles[fd] = f
		res.Return = fd
	} else {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	}

	return res
}
