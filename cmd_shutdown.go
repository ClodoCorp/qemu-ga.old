package main

import (
	"encoding/json"
	"os/exec"
)

var cmdShutdown = &Command{
	Name: "guest-shutdown",
	Func: fnShutdown,
}

func init() {
	commands = append(commands, cmdShutdown)
}

func fnShutdown(d map[string]interface{}) interface{} {
	id, _ := (d["id"].(json.Number)).Int64()
	mode := d["mode"].(string)

	var args []string = []string{"-h"}

	switch mode {
	case "halt":
		args = append(args, "-H")
		break
	case "reboot":
		args = append(args, "-r")
		break
	case "powerdown":
	default:
		args = append(args, "-P")
		break
	}
	args = append(args, "+0", "hypervisor initiated shutdown")
	cmd := exec.Command("shutdown", args...)
	defer cmd.Run()

	return &Response{
		Return: id,
	}
}
