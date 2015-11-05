/*
guest-shutdown - shutdown guest via agent

Example:
        { "execute": "guest-shutdown", "arguments": {
            "mode": string // optional, shutdown mode (halt, reboot, powerdown), default powerdown
          }
        }
*/
package guest_shutdown

import (
	"encoding/json"
	"os/exec"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-shutdown",
		Func:    fnGuestShutdown,
		Enabled: true,
	})
}

func fnGuestShutdown(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Mode string `json:"mode"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{"-h"}

	switch reqData.Mode {
	case "halt":
		args = append(args, "-H")
		break
	case "reboot":
		args = append(args, "-r")
		break
	case "powerdown":
	default:
		args = append(args, "-P")
		break
	}
	args = append(args, "+0", "hypervisor initiated shutdown")
	cmd := exec.Command("shutdown", args...)
	defer cmd.Run()

	return res
}
