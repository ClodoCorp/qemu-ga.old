/*
guest-stats - returns disk and memory stats from guest

Example:
        { "execute": "guest-stats", "arguments": {}}
*/
package guest_stats

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-stats",
		Func:    fnGuestStats,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestStats(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}
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
		La1         float64
		La5         float64
		La15        float64
	}{}

	buf, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	fields := strings.Fields(string(buf))
	if resData.La1, err = strconv.ParseFloat(fields[0], 64); err != nil {
		resData.La1 = float64(-1)
	}
	if resData.La5, err = strconv.ParseFloat(fields[1], 64); err != nil {
		resData.La5 = float64(-1)
	}
	if resData.La15, err = strconv.ParseFloat(fields[2], 64); err != nil {
		resData.La15 = float64(-1)
	}

	buf, err = ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
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
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	resData.BlkTotal = st.Blocks * uint64(st.Frsize)
	resData.BlkFree = st.Bavail * uint64(st.Frsize)

	resData.InodeTotal = st.Files
	resData.InodeFree = st.Ffree

	res.Return = resData
	return res
}
