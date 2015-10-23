package main

import (
	"encoding/json"
	"fmt"
)

var cmdFileClose = &Command{
	Name:    "guest-file-close",
	Func:    fnFileClose,
	Enabled: true,
}

func init() {
	commands = append(commands, cmdFileClose)
}

func fnFileClose(req *Request) *Response {
	res := &Response{Id: req.Id}

	file := struct {
		Handle int `json:"handle"`
	}{}

	err := json.Unmarshal(req.RawArgs, &file)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[file.Handle]; ok {
			if err = f.Close(); err != nil {
				res.Error = &Error{Code: -1, Desc: err.Error()}
			} else {
				delete(openFiles, file.Handle)
			}
		} else {
			res.Error = &Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
