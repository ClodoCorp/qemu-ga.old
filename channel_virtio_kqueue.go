// +build freebsd netbsd openbsd
// +build !linux

package main

import "fmt"

func (ch *VirtioChannel) Poll() error {
	return fmt.Errorf("not implemented")
}
