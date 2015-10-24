package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

var cmdFileRead = &Command{
	Name:    "guest-file-read",
	Func:    fnFileRead,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdFileRead)
}

func fnFileRead(req *Request) *Response {
	res := &Response{Id: req.Id}

	file := struct {
		Handle int `json:"handle"`
		Count  int `json:"count,omitempty"`
	}{}

	ret := struct {
		Count  int    `json:"count"`
		BufB64 string `json:"buf-b64"`
		Eof    bool   `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &file)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[file.Handle]; ok {
			var buffer []byte
			var n int
			n, err = f.Read(buffer)
			switch err {
			case nil:
				ret.Count = n
				ret.BufB64 = base64.StdEncoding.EncodeToString(buffer)
				res.Return = ret
			case io.EOF:
				ret.Count = n
				ret.BufB64 = base64.StdEncoding.EncodeToString(buffer)
				ret.Eof = true
				res.Return = ret
			default:
				res.Error = &Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
