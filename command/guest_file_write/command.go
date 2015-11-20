/*
guest-file-write - write file to fd inside guest

Example:
        { "execute": "guest-file-write", "arguments": {
            "handle": int // required, unique fd identifier
            "buf-b64": string // required, base64 encoded data
            "count": int // optional, number of bytes to write
          }
        }
*/
package guest_file_write

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-file-write",
		Func:      fnGuestFileWrite,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestFileWrite(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Handle int    `json:"handle"`
		BufB64 string `json:"buf-b64"`
		Count  int    `json:"count,omitempty"`
	}{}

	resData := struct {
		Count int  `json:"count"`
		Eof   bool `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	} else {
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			var buffer []byte
			buffer, err = base64.StdEncoding.DecodeString(reqData.BufB64)
			if err != nil {
				res.Error = &qga.Error{Code: -1, Desc: err.Error()}
				return res
			}
			var n int
			n, err = f.Write(buffer)
			switch err {
			case nil:
				resData.Count = n
				res.Return = resData
			case io.EOF:
				resData.Count = n
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
