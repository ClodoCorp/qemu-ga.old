package main

import (
	"encoding/json"
	"os/exec"
)

var cmdFstrim = &Command{
	Name:    "guest-fstrim",
	Func:    fnFstrim,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdFstrim)
}

// TODO: USE NATIVE SYSCALL
func fnFstrim(req *Request) *Response {
	res := &Response{Id: req.Id}
	//	r := ioctl.FsTrimRange{Start: 0, Length: -1, MinLength: 0}

	reqData := struct {
		Minimum int `json:"minimum,omitempty"`
	}{}

	type Path struct {
		Path    string `json:"path"`
		Trimmed *int   `json:"trimmed,omitempty"`
		Minimum *int   `json:"minimum,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	resData := struct {
		Paths []*Path `json:"paths"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	fslist, err := listMountedFileSystems()
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	/*
		if f, err := os.OpenFile("/", os.O_RDONLY, os.FileMode(0400)); err == nil {
			defer f.Close()
			err = ioctl.Fitrim(uintptr(f.Fd()), uintptr(unsafe.Pointer(&r)))
	*/
	for _, fs := range fslist {
		switch fs.Type {
		case "ufs", "ffs":
			err = exec.Command("fsck_"+fs.Type, "-B", "-E", fs.Path).Run()
		default:
			err = exec.Command("fstrim", fs.Path).Run()
		}
		rpath := &Path{Path: fs.Path}
		if err != nil {
			rpath.Error = err.Error()
		}
		resData.Paths = append(resData.Paths, rpath)
	}

	res.Return = resData
	return res
}
