package main

var cmdNetRoute = &Command{
	Name: "guest-network-set-routes",
	Func: fnNetRoute,
}

func init() {
	commands = append(commands, cmdNetRoute)
}

func fnNetRoute(d map[string]interface{}) interface{} {
	return &Response{}
}
