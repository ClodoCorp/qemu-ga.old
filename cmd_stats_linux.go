package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
)

var cmdStats = &Command{
	Name: "guest-stats",
	Func: fnStats,
}

func init() {
	commands = append(commands, cmdStats)
}

func fnStats(d map[string]interface{}) interface{} {
	var st syscall.Statfs_t
	type StatsInfo struct {
		MemoryTotal uint64
		MemoryFree  uint64
		SwapTotal   uint64
		SwapFree    uint64
		BlkTotal    uint64
		BlkFree     uint64
		InodeTotal  uint64
		InodeFree   uint64
	}
	res := &StatsInfo{}
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
		value, err := strconv.ParseUint(strings.TrimSpace(fields[1]), 10, 64)
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

	err = syscall.Statfs("/", &st)
	if err != nil {
		return &Response{}
	}

	res.BlkTotal = st.Blocks * uint64(st.Frsize)
	res.BlkFree = st.Bavail * uint64(st.Frsize)

	res.InodeTotal = st.Files
	res.InodeFree = st.Ffree

	return &Response{Return: res}
}
