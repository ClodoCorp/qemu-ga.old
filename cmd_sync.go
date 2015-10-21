package main

import "encoding/json"

var cmdSync = &Command{
	Name: "guest-sync",
	Func: fnSync,
}

func init() {
	commands = append(commands, cmdSync)
}

func fnSync(req *Request) *Response {
	res := &Response{}

	sync := struct {
		Id int `json:"id"`
	}{}

	err := json.Unmarshal(req.RawArgs, &sync)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
	} else {
		res.Return = sync.Id
		res.Id = req.Id
	}

	return res
}
