package main

var cmdPing = &Command{
	Name:    "guest-ping",
	Func:    fnPing,
	Enabled: true,
}

func init() {
	commands = append(commands, cmdPing)
}

func fnPing(req *Request) *Response {
	return &Response{}
}
