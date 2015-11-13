package main

import (
	"encoding/json"
	"os/exec"

	"golang.org/x/sys/unix"
)

var cmdSetTime = &Command{
	Name:    "guest-set-time",
	Func:    fnSetTime,
	Enabled: true,
}

func init() {
	commands = append(commands, cmdSetTime)
}

func fnSetTime(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := reqDataSetTime{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{}
	if reqData.Time != 0 {
		tv := &unix.Timeval{Sec: reqData.Time / 1000000000, Usec: reqData.Time % 1000000000 / 1000}
		if err = unix.Settimeofday(tv); err != nil {
			res.Error = &Error{Code: -1, Desc: err.Error()}
			return res
		}
		args = append(args, "-w")
	} else {
		args = append(args, "-s")
	}

	err = exec.Command("hwclock", args...).Run()
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	return res
}
