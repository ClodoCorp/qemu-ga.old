package main

import (
	"encoding/json"
	"os"
	"unsafe"

	"github.com/vtolstov/go-ioctl"
)

var cmdFstrim = &Command{
	Name: "guest-fstrim",
	Func: fnFstrim,
}

func init() {
	commands = append(commands, cmdFstrim)
}

func fnFstrim(d map[string]interface{}) interface{} {
	var r int

	id, _ := (d["id"].(json.Number)).Int64()

	if f, err := os.Open("/dev/sda1"); err == nil {
		defer f.Close()
		ioctl.Fitrim(f.Fd(), uintptr(unsafe.Pointer(&r)))
	}

	return &Response{
		Return: id,
	}
}
