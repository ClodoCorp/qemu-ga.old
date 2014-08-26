package main

import "encoding/json"

var cmdSync = &Command{
	Name: "guest-sync",
	Func: fnSync,
}

func init() {
	commands = append(commands, cmdSync)
}

func fnSync(d map[string]interface{}) interface{} {
	id, _ := (d["id"].(json.Number)).Int64()
	return &Response{
		Return: id,
	}
}
