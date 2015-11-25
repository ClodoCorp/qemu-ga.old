// +build linux freebsd netbsd openbsd windows

package guest_set_time

type reqDataSetTime struct {
	Time int64 `json:"time,omitempty"`
}
