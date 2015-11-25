// +build linux freebsd

package guest_set_time

import "golang.org/x/sys/unix"

func newTimeval(t int64) *unix.Timeval {
	return &unix.Timeval{Sec: t / 1000000000, Usec: t % 1000000000 / 1000}
}
