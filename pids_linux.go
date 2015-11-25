package main

import (
	"os"

	"github.com/vtolstov/go-ps"
)

func getPids(name string, filter bool) []int {
	pids := []int{}

	procs, err := ps.FindProcessByExecutable(name)
	if err != nil || len(procs) == 0 {
		return pids
	}

	ownpid := os.Getpid()

Check:
	for _, proc := range procs {
		if filter {
			for _, pid := range proc.CPids() {
				if pid == ownpid {
					continue Check
				}
			}
		}
		pids = append(pids, proc.Pid())
		pids = append(pids, proc.CPids()...)
	}

	return pids
}
