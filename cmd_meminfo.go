package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var cmdMemInfo = &Command{
	Name: "guest-memory-info",
	Func: fnMemInfo,
}

func init() {
	commands = append(commands, cmdMemInfo)
}

func fnMemInfo(d map[string]interface{}) interface{} {
	type MemoryInfo struct {
		MemoryTotal int64
		MemoryFree  int64
		SwapTotal   int64
		SwapFree    int64
	}
	res := &MemoryInfo{}
	buf, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return &Response{}
	}

	reader := bufio.NewReader(bytes.NewBuffer(buf))

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		fields := strings.Fields(string(line))
		fmt.Printf("%+v\n", fields)
		value, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 64)
		if err != nil {
			fmt.Printf("err %s\n", err)
		}
		switch strings.TrimSpace(fields[0]) {
		case "MemTotal:":
			res.MemoryTotal = value * 1024
		case "MemFree:", "Cached:", "Buffers:":
			res.MemoryFree += value * 1024
		case "SwapTotal:":
			res.SwapTotal = value * 1024
		case "SwapFree:":
			res.SwapFree = value * 1024
		}
	}

	return &Response{Return: res}
}
