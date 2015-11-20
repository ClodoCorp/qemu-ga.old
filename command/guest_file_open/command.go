/*
guest-file-open - open file inside guest and returns it handle

Example:
        { "execute": "guest-file-open", "arguments": {
            "path": string // required, file path
            "mode": string // optional, file open mode
          }
        }
*/
package guest_file_open

import (
	"encoding/json"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-open",
		Func:      fnGuestFileOpen,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestFileOpen(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Path string `json:"path"`
		Mode string `json:"mode,omitempty"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	var flag int
	if reqData.Mode != "" {
		for _, s := range reqData.Mode {
			switch s {
			case 'a':
				flag = flag | os.O_APPEND | os.O_CREATE | os.O_WRONLY
			case '+':
				flag = flag | os.O_RDWR
			case 'w':
				flag = flag | os.O_TRUNC | os.O_WRONLY
			case 'r':
				flag = flag | os.O_RDONLY
			}
		}
	} else {
		flag = flag | os.O_RDONLY
	}

	if f, err := os.OpenFile(reqData.Path, flag, os.FileMode(0600)); err == nil {
		fd := int(f.Fd())
		qga.StoreSet("guest-file", fd, f)
		res.Return = fd
	} else {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	}

	return res
}
