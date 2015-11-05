/*
guest-get-time - get current guest time in nanoseconds

Example:
        { "execute": "guest-get-time", "arguments": {}}
*/
package guest_set_time

import (
	"time"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-get-time",
		Func:    fnGuestGetTime,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestGetTime(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	res.Return = struct {
		Time int64
	}{Time: time.Now().UnixNano()}

	return res
}
