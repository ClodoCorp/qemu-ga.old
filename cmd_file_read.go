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

	reqData := struct {
		Handle int `json:"handle"`
		Count  int `json:"count,omitempty"`
	}{}

	resData := struct {
		Count  int    `json:"count"`
		BufB64 string `json:"buf-b64"`
		Eof    bool   `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[reqData.Handle]; ok {
			var buffer []byte
			var n int
			n, err = f.Read(buffer)
			switch err {
			case nil:
				resData.Count = n
				resData.BufB64 = base64.StdEncoding.EncodeToString(buffer)
				res.Return = resData
			case io.EOF:
				resData.Count = n
				resData.BufB64 = base64.StdEncoding.EncodeToString(buffer)
				resData.Eof = true
				res.Return = resData
			default:
				res.Error = &Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
