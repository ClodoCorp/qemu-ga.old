// +build ignore

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type Dev struct {
	Major uint64
	Minor uint64
}

func (d *Dev) String() string {
	return fmt.Sprintf("%d:%d", d.Major, d.Minor)
}

func (d *Dev) Int() int {
	return int(d.Major*256 + d.Minor)
}

func main() {
	var devFs *Dev
	var devBlk *Dev
	var err error
	devFs, err = findFs()
	if err != nil {
		panic(err)
	}
	devBlk, err = findBlock("/sys/block", devFs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s\n", devFs, devBlk)
	err = syscall.Mknod("/tmp/block", uint32(os.ModeDevice|syscall.S_IFBLK|0600), devBlk.Int())
	if err != nil {
		panic(err)
	}
}

func findFs() (*Dev, error) {
	var st syscall.Stat_t

	err := syscall.Stat("/", &st)
	if err != nil {
		return nil, err
	}
	return &Dev{Major: uint64(st.Dev / 256), Minor: uint64(st.Dev % 256)}, nil
}

func findBlock(start string, s *Dev) (*Dev, error) {
	var err error
	fis, err := ioutil.ReadDir(start)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		switch fi.Name() {
		case "bdi", "subsystem", "device", "trace":
			continue
		}
		if _, err := os.Stat(filepath.Join(start, "dev")); err == nil {
			if buf, err := ioutil.ReadFile(filepath.Join(start, "dev")); err == nil {
				dev := strings.TrimSpace(string(buf))
				if s.String() == dev {
					if buf, err = ioutil.ReadFile(filepath.Join(filepath.Dir(start), "dev")); err == nil {
						majorminor := strings.Split(strings.TrimSpace(string(buf)), ":")
						major, _ := strconv.Atoi(majorminor[0])
						minor, _ := strconv.Atoi(majorminor[1])
						return &Dev{Major: uint64(major), Minor: uint64(minor)}, nil
					}
				}
			}
		}
		devBlk, err := findBlock(filepath.Join(start, fi.Name()), s)
		if err == nil {
			return devBlk, err
		}
	}
	return nil, errors.New("failed to find dev")
}
