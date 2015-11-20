/*
guest-file-chown - set file owner

Example:
        { "execute": "guest-file-chown", "arguments": {
            "uid": int // required, file owner uid
            "gid": int // required, file owner gid
            "path": string // optional, file path
            "handle": int // optional, file handle
          }
        }
*/
package guest_file_chown

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-chown",
		Func:      fnGuestFileChown,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestFileChown(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Uid    int    `json:"uid"`
		Gid    int    `json:"gid"`
		Handle int    `json:"handle,omitempty"`
		Path   string `json:"path,omitempty"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	if reqData.Path == "" && reqData.Handle == 0 {
		res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("either path or handle must be non empty")}
		return res
	}

	switch {
	case reqData.Path != "":
		if err = os.Chown(reqData.Path, reqData.Uid, reqData.Gid); err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		}
	case reqData.Handle != 0:
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			if err = f.Chown(reqData.Uid, reqData.Gid); err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
