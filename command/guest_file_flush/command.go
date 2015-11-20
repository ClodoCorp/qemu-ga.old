/*
guest-file-flush - flush buffered data to file

Example:
        { "execute": "guest-file-flush", "arguments": {
            "handle": int // required, unique fd identifier
          }
        }
*/
package guest_file_flush

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-flush",
		Func:      fnGuestFileFlush,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestFileFlush(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Handle int `json:"handle"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	} else {
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			if err = f.Sync(); err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
