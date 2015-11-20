/*
guest-file-chmod - set file mode

Example:
        { "execute": "guest-file-chmod", "arguments": {
            "mode": int // required, file mode
            "path": string // optional, file path
            "handle": int // optional, file handle
          }
        }
*/
package guest_file_chmod

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-chmod",
		Func:      fnGuestFileChmod,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestFileChmod(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Mode   int    `json:"mode"`
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
		if err = os.Chmod(reqData.Path, os.FileMode(reqData.Mode)); err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		}
	case reqData.Handle != 0:
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			if err = f.Chmod(os.FileMode(reqData.Mode)); err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
