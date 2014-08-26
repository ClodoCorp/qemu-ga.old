package main

var cmdPing = &Command{
	Name: "guest-ping",
	Func: fnPing,
}

func init() {
	commands = append(commands, cmdPing)
}

func fnPing(d map[string]interface{}) interface{} {
	return &Response{}
}
