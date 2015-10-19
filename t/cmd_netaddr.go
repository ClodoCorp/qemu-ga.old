package main

var cmdNetAddr = &Command{
	Name: "guest-network-set-addresses",
	Func: fnNetAddr,
}

func init() {
	commands = append(commands, cmdNetAddr)
}

func fnNetAddr(d map[string]interface{}) interface{} {
	return &Response{}
}
