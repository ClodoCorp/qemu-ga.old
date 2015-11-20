/*
guest-file-read - read file inside guest via fd

Example:
        { "execute": "guest-file-read", "arguments": {
            "handle": int // required, unique fd identifier
            "count": int // optional, bytes count to read
          }
        }
*/
package guest_file_read

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
		Name:      "guest-file-read",
		Func:      fnGuestFileRead,
		Enabled:   true,
		Returns:   true,
		Arguments: true,
	})
}

func fnGuestFileRead(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		Handle int `json:"handle"`
		Count  int `json:"count,omitempty"`
	}{}

	resData := struct {
		Count  int    `json:"count"`
		BufB64 string `json:"buf-b64"`
		Eof    bool   `json:"eof"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
	} else {
		if iface, ok := qga.StoreGet("guest-file", reqData.Handle); ok {
			f := iface.(*os.File)
			var buffer []byte
			var n int
			n, err = f.Read(buffer)
			switch err {
			case nil:
				resData.Count = n
				resData.BufB64 = base64.StdEncoding.EncodeToString(buffer)
				res.Return = resData
			case io.EOF:
				resData.Count = n
				resData.BufB64 = base64.StdEncoding.EncodeToString(buffer)
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
