package main

import (
	"encoding/json"
	"fmt"
)

var cmdFileFlush = &Command{
	Name:    "guest-file-flush",
	Func:    fnFileFlush,
	Enabled: true,
}

func init() {
	commands = append(commands, cmdFileFlush)
}

func fnFileFlush(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		Handle int `json:"handle"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[reqData.Handle]; ok {
			if err = f.Sync(); err != nil {
				res.Error = &Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
