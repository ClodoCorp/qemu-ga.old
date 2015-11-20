/*
guest-file-close - close file handle

Example:
        { "execute": "guest-file-close", "arguments": {
            "handle": int // required, unique fd identifier
          }
        }
*/
package guest_file_close

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-close",
		Func:      fnGuestFileClose,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestFileClose(req *qga.Request) *qga.Response {
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
			if err = f.Close(); err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			} else {
				qga.StoreDel("guest-file", reqData.Handle)
			}
		} else {
			res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
