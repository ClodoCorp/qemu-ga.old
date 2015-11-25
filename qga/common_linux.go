// +build linux
// +build !freebsd !netbsd !openbsd

package qga

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ListMountedFileSystems() ([]FileSystem, error) {
	var fs []FileSystem
	var line string

	f, err := os.Open("/proc/self/mounts")
	if err != nil {
		return fs, err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		if line, err = br.ReadString('\n'); err != nil {
			break
		}

		values := strings.Fields(line)
		if values[1] != "/" {
			continue
		}
		switch values[2] {
		case "tmpfs", "cgroup", "debugfs", "smbfs", "cifs", "rootfs":
			continue

		}
		/*
		   Device  string
		   Path    string
		   Type    string
		   Options []string
		*/
		fs = append(fs, FileSystem{Device: values[0], Path: values[1], Type: values[2], Options: strings.Split(values[3], ",")})

	}
	if err == io.EOF {
		err = nil
	}
	return fs, err
}
