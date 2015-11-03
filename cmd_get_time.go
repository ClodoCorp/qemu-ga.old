package main

import "time"

var cmdGetTime = &Command{
	Name:    "guest-get-time",
	Func:    fnGetTime,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdGetTime)
}

func fnGetTime(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		ID int64 `json:"id"`
	}{}

	resData := struct {
		Time int64
	}{Time: time.Now().UnixNano()}

	res.Return = resData.Time

	return res
}
