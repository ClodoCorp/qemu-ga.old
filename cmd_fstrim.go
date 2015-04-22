package main

import (
	"encoding/json"
	"os/exec"
)

var cmdFstrim = &Command{
	Name: "guest-fstrim",
	Func: fnFstrim,
}

func init() {
	commands = append(commands, cmdFstrim)
}

// TODO: USE NATIVE SYSCALL
func fnFstrim(d map[string]interface{}) interface{} {
	//	r := ioctl.FsTrimRange{Start: 0, Length: -1, MinLength: 0}
	id, _ := (d["id"].(json.Number)).Int64()
	fslist, err := listMountedFileSystems()
	if err != nil {
		return &Response{}
	}
	/*
		if f, err := os.OpenFile("/", os.O_RDONLY, os.FileMode(0400)); err == nil {
			defer f.Close()
			err = ioctl.Fitrim(uintptr(f.Fd()), uintptr(unsafe.Pointer(&r)))
	*/
	for _, fs := range fslist {
		switch fs.Type {
		case "ufs", "ffs":
			exec.Command("fsck_"+fs.Type, "-B", "-E", fs.Path).Run()
		default:
			exec.Command("fstrim", fs.Path).Run()
		}
	}
	return &Response{
		Return: id,
	}
}
