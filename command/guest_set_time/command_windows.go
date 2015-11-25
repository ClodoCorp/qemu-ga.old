/*
guest-set-time - set guest time

Example:
        { "execute": "guest-set-time", "arguments": {
            "time": int // optional, time to set
          }
        }
*/
package guest_set_time

import "github.com/vtolstov/qemu-ga/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-set-time",
		Func:    fnGuestSetTime,
		Enabled: false,
	})
}

func fnGuestSetTime(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	return res
}
