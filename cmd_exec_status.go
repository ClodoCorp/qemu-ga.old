package main

import "encoding/json"

var cmdExecStatus = &Command{
	Name:    "guest-exec-status",
	Func:    fnExecStatus,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdExecStatus)
}

func fnExecStatus(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		Pid int `json:"pid"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	if s, ok := execStatuses[reqData.Pid]; ok {
		res.Return = s
		if s.Exited {
			delete(execStatuses, reqData.Pid)
		}
	} else {
		res.Error = &Error{Code: -1, Desc: "provided pid not found"}
	}

	return res
}
