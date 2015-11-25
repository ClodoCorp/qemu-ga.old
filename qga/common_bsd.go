// +build netbsd freebsd openbsd
// +build !linux

package qga

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

func ListMountedFileSystems() ([]FileSystem, error) {
	var fs []FileSystem
	var line string

	cmd := exec.Command("mount", "-p")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fs, err
	}
	br := bufio.NewReader(stdout)
	if err = cmd.Start(); err != nil {
		return fs, err
	}

	for {
		if line, err = br.ReadString('\n'); err != nil {
			break
		}

		values := strings.Fields(line)

		if values[0] != "/" {
			continue
		}

		switch values[2] {
		case "tmpfs", "devfs":
			continue
		}

		fs = append(fs, FileSystem{Device: values[0], Path: values[1], Type: values[2], Options: strings.Split(values[3], ",")})

	}

	if err == io.EOF {
		err = nil
	}

	if err = cmd.Wait(); err != nil {
		return fs, err
	}

	return fs, err
}
