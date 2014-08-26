package main

import (
	"os"
)

func init() {
	if f, e := os.OpenFile("/proc/self/oom_score_adj", os.O_WRONLY, 0644); e == nil {
		f.Write([]byte("-1000"))
		f.Close()
		return
	}
	if f, e := os.OpenFile("/proc/self/oom_adj", os.O_WRONLY, 0644); e == nil {
		f.Write([]byte("-17"))
		f.Close()
		return
	}
}
