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
	res := &Response{}
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
	stinfo := &StatsInfo{}
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
			stinfo.MemoryTotal = value * 1024
		case "MemFree:", "Cached:", "Buffers:":
			stinfo.MemoryFree += value * 1024
		case "SwapTotal:":
			stinfo.SwapTotal = value * 1024
		case "SwapFree:":
			stinfo.SwapFree = value * 1024
		}
	}

	err = syscall.Statfs("/", &st)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	stinfo.BlkTotal = st.Blocks * uint64(st.Frsize)
	stinfo.BlkFree = st.Bavail * uint64(st.Frsize)

	stinfo.InodeTotal = st.Files
	stinfo.InodeFree = st.Ffree

	res.Return = stinfo
	res.Id = req.Id
	return res
}
