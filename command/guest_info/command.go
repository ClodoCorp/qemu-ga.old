/*
guest-info - request agent info from guest

Example:
        { "execute": "guest-info", "arguments": {}}
*/
package guest_info

import (
	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-info",
		Func:    fnGuestInfo,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestInfo(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	res.Return = struct {
		Version  string         `json:"version"`
		Commands []*qga.Command `json:"supported_commands"`
	}{Version: qga.GetVersion(), Commands: qga.ListCommands()}

	return res
}
