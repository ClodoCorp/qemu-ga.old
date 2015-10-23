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

	file := struct {
		Handle int `json:"handle"`
		Offset int `json:"offset"`
		Whence int `json:"whence"`
	}{}

	ret := struct {
		Pos int  `json:"position"`
		Eof bool `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &file)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		if f, ok := openFiles[file.Handle]; ok {
			n, err := f.Seek(int64(file.Offset), file.Whence)
			switch err {
			case nil:
				ret.Pos = int(n)
				res.Return = ret
			case io.EOF:
				ret.Pos = int(n)
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
