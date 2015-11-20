/*
guest-exec-status - get status from running command

Example:
        { "execute": "guest-exec-status", "arguments": {
            "pid": int // required, process id from guest-exec
          }
        }
*/
package guest_exec_status

import (
	"encoding/json"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-exec-status",
		Func:      fnGuestExecStatus,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestExecStatus(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Pid int `json:"pid"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	if iface, ok := qga.StoreGet("guest-exec", reqData.Pid); ok {
		s := iface.(*qga.ExecStatus)
		res.Return = s
		if s.Exited {
			qga.StoreDel("guest-exec", reqData.Pid)
		}
	} else {
		res.Error = &qga.Error{Code: -1, Desc: "provided pid not found"}
	}

	return res
}
