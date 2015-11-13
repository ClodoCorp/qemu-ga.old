package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

var cmdMemInfo = &Command{
	Name:    "guest-memory-info",
	Func:    fnMemInfo,
	Enabled: true,
	Returns: true,
}

func init() {
	commands = append(commands, cmdMemInfo)
}

func fnMemInfo(req *Request) *Response {
	res := &Response{Id: req.Id}

	resData := struct {
		MemoryTotal int64
		MemoryFree  int64
		SwapTotal   int64
		SwapFree    int64
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
		value, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64)
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

	res.Return = resData
	return res
}
