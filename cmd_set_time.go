package main

import (
	"encoding/json"

	"golang.org/x/sys/unix"
)

var cmdSetTime = &Command{
	Name:    "guest-set-time",
	Func:    fnSetTime,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdSetTime)
}

func fnSetTime(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		Time int64 `json:"time,omitempty"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	if req.Data.Time != 0 {
		tv := &unix.Timeval{Sec: req.Data.Time / 1000000000, Usec: (req.Data.Time % 1000000000) / 1000}
		if err = unix.Settimeofday(tv); err != nil {
			res.Error = &Error{Code: -1, Desc: err.Error()}
			return res
		}
	}

	res.Return = resData.Time

	return res
}
