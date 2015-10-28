package main

import (
	"encoding/json"
	"fmt"
	"io"
)

var cmdFileSeek = &Command{
	Name:    "guest-file-seek",
	Func:    fnFileSeek,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdFileSeek)
}

func fnFileSeek(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		Handle int `json:"handle"`
		Offset int `json:"offset"`
		Whence int `json:"whence"`
	}{}

	resData := struct {
		Pos int  `json:"position"`
		Eof bool `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[reqData.Handle]; ok {
			n, err := f.Seek(int64(reqData.Offset), reqData.Whence)
			switch err {
			case nil:
				resData.Pos = int(n)
				res.Return = resData
			case io.EOF:
				resData.Pos = int(n)
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
