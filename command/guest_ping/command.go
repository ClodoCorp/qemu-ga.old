package guest_ping

import "github.com/vtolstov/qemu-ga/qga"

func init() {
	qga.RegisterCommand(&Command{
		Name:    "guest-ping",
		Func:    fnGuestPing,
		Enabled: true,
	})
}

func fnGuestPing(req *qga.Request) *qga.Response {
	return &qga.Response{}
}
