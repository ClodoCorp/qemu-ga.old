/*
guest-file-seek - seek on file

Example:
        { "execute": "guest-file-seek", "arguments": {
            "handle": int // required, unique fd identifier
            "offset": int // required, offset inside file
            "whence": int // required, starting point to seek (-1, 0, 1)
          }
        }
*/
package guest_file_seek

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-seek",
		Func:      fnGuestFileSeek,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestFileSeek(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Handle int `json:"handle"`
		Offset int `json:"offset"`
		Whence int `json:"whence"`
	}{}

	resData := struct {
		Pos int  `json:"position"`
		Eof bool `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	} else {
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			n, err := f.Seek(int64(reqData.Offset), reqData.Whence)
			switch err {
			case nil:
				resData.Pos = int(n)
				res.Return = resData
			case io.EOF:
				resData.Pos = int(n)
				resData.Eof = true
				res.Return = resData
			default:
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			}
		} else {
			res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("file handle not found")}
		}
	}

	return res
}
