/*
guest-sync - sync host<->guest communication

Example:
        { "execute": "guest-sync", "arguments": {
            "id": int // required, unique identifier
          }
        }
*/
package guest_sync

import (
	"encoding/json"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-sync",
		Func:      fnGuestSync,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestSync(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		ID int64 `json:"id"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	} else {
		res.Return = reqData.ID
	}

	return res
}
