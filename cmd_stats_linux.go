package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
)

var cmdStats = &Command{
	Name:    "guest-stats",
	Func:    fnStats,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdStats)
}

func fnStats(req *Request) *Response {
	res := &Response{Id: req.Id}
	var st syscall.Statfs_t

	resData := struct {
		MemoryTotal uint64
		MemoryFree  uint64
		SwapTotal   uint64
		SwapFree    uint64
		BlkTotal    uint64
		BlkFree     uint64
		InodeTotal  uint64
		InodeFree   uint64
	}{}

	buf, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	reader := bufio.NewReader(bytes.NewBuffer(buf))

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		fields := strings.Fields(string(line))
		value, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
		if err != nil {
			continue
		}
		switch strings.TrimSpace(fields[0]) {
		case "MemTotal:":
			resData.MemoryTotal = value * 1024
		case "MemFree:", "Cached:", "Buffers:":
			resData.MemoryFree += value * 1024
		case "SwapTotal:":
			resData.SwapTotal = value * 1024
		case "SwapFree:":
			resData.SwapFree = value * 1024
		}
	}

	err = syscall.Statfs("/", &st)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	resData.BlkTotal = st.Blocks * uint64(st.Frsize)
	resData.BlkFree = st.Bavail * uint64(st.Frsize)

	resData.InodeTotal = st.Files
	resData.InodeFree = st.Ffree

	res.Return = resData
	return res
}
