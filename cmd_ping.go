package main

var cmdPing = &Command{
	Name: "guest-ping",
	Func: fnPing,
}

func init() {
	commands = append(commands, cmdPing)
}

func fnPing(req *Request) *Response {
	ret := struct {
		id int `json:"-"`
	}{}
	return &Response{Return: ret}
}
