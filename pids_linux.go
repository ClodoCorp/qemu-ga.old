package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getPids(pattern string, filter bool) (pids []int) {
	pid := os.Getpid()

	fis, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Printf("err %s\n", err.Error())
		return pids
	}
Check:
	for _, fi := range fis {
		chpid, err := strconv.Atoi(fi.Name())
		if err != nil {
			continue Check
		}

		buf, err := ioutil.ReadFile(filepath.Join("/proc", fi.Name(), "comm"))

		if err == nil {
			progname := strings.TrimSpace(string(buf))
			switch pattern {
			case progname, filepath.Base(progname):
				if filter {
					ffis, err := ioutil.ReadDir(filepath.Join("/proc", fi.Name(), "task"))
					if err == nil {
						for _, ffi := range ffis {
							if ffi.Name() == fmt.Sprintf("%d", pid) {
								continue Check
							}
						}
						for _, ffi := range ffis {
							if ffi.Name() != fi.Name() {
								if cchpid, err := strconv.Atoi(ffi.Name()); err == nil {
									pids = append(pids, cchpid)
								}
							}
						}
					}
				}
				pids = append(pids, chpid)
			}

		}
	}
	return pids
}
