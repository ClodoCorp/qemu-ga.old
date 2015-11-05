package main

import "encoding/json"

var cmdSync = &Command{
	Name:    "guest-sync",
	Func:    fnSync,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdSync)
}

func fnSync(req *Request) *Response {
	res := &Response{Id: req.Id}

	reqData := struct {
		ID int64 `json:"id"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		res.Return = reqData.ID
	}

	return res
}
