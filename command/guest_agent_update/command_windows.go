/*
guest-agent-update - update qemu-ga inside vm

Example:
        { "execute": "guest-agent-update", "arguments": {
            "path": string // required, http/https/file path to qemu-ga binary for update
            "timeout": int // optional, timeout for http/https transport
          }
        }
*/
package guest_agent_update

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-agent-update",
		Func:      fnGuestAgentUpdate,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestAgentUpdate(req *qga.Request) *qga.Response {
	res := &qga.Response{}
	var r io.ReadCloser
	var httpClient *http.Client

	httpTransport := &http.Transport{
		Dial:            (&net.Dialer{DualStack: true}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	reqData := struct {
		Path    string `json:"path"`
		Timeout int64  `json:"timeout,omitempty"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	if reqData.Timeout == 0 {
		reqData.Timeout = 30
	}

	dt, err := time.ParseDuration(fmt.Sprintf("%ds", reqData.Timeout))
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	httpClient = &http.Client{Transport: httpTransport, Timeout: dt}

	u, err := url.Parse(reqData.Path)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	switch u.Scheme {
	case "http", "https":
		hres, err := httpClient.Get(reqData.Path)
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		r = hres.Body
	case "file":
		r, err = os.Open(u.Path)
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
	default:
		res.Error = &qga.Error{Code: -1, Desc: fmt.Sprintf("invalid path %s", u)}
		return res
	}
	defer r.Close()

	dirname, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	filename := fmt.Sprintf(".%s", filepath.Base(os.Args[0]))
	w, err := os.OpenFile(filepath.Join(dirname, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0755))
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	_, err = io.Copy(w, r)
	if err != nil {
		defer w.Close()
		defer os.Remove(filepath.Join(dirname, filename))
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	w.Sync()
	w.Close()

	if err = os.Rename(filepath.Join(dirname, filename), filepath.Join(dirname, filepath.Base(os.Args[0]))); err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	time.Sleep(2 * time.Second)
	defer func() {
		cmd := exec.Command(os.Args[0])

		err = cmd.Start()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()

	return res
}
