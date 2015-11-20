/*
guest-file-stat - get stat on file

Example:
        { "execute": "guest-file-stat", "arguments": {
            "path": string // optional, file path
            "handle": int // optional, file handle
          }
        }
*/
package guest_file_stat

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-stat",
		Func:      fnGuestFileStat,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestFileStat(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}
	var fi os.FileInfo

	reqData := struct {
		Handle int    `json:"handle,omitempty"`
		Path   string `json:"path,omitempty"`
	}{}

	resData := struct {
		Name   string `json:"name"`
		Size   int64  `json:"size"`
		Mode   uint32 `json:"mode"`
		Uid    int    `json:"uid,omitempty"`
		Gid    int    `json:"gid,omitempty"`
		Modify int64  `json:"modify"`
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
		if fi, err = os.Stat(reqData.Path); err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
	case reqData.Handle != 0:
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			if fi, err = f.Stat(); err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
				return res
			}
		} else {
			res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
			return res
		}
	}

	resData.Name = fi.Name()
	resData.Size = fi.Size()
	resData.Mode = uint32(fi.Mode())
	resData.Modify = fi.ModTime().Unix()
	res.Return = resData
	return res
}
