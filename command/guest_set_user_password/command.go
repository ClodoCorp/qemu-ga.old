/*
guest-set-user-password - sync host<->guest communication

Example:
        { "execute": "guest-set-user-password", "arguments": {
            "username": string // required, username to change password
            "password": string // required, base64 encoded password
            "crypted": bool // optional, specify that password already encrypted
          }
        }
*/
package guest_set_user_password

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/vtolstov/qemu-ga/qga"
)

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:      "guest-set-user-password",
		Func:      fnGuestSetUserPassword,
		Enabled:   true,
		Arguments: true,
	})
}

func fnGuestSetUserPassword(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}

	reqData := struct {
		User    string `json:"username"`
		Passwd  string `json:"password"`
		Crypted bool   `json:"crypted"`
	}{}

	err := json.Unmarshal(req.RawArgs, &reqData)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	passwd, err := base64.StdEncoding.DecodeString(reqData.Passwd)
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	args := []string{}

	if reqData.Crypted {
		args = append(args, "-e")
	}

	cmd := exec.Command("chpasswd", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	err = cmd.Start()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	arg := fmt.Sprintf("%s:%s", reqData.User, passwd)
	_, err = stdin.Write([]byte(arg))
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		res.Error = &qga.Error{Code: -1, Desc: err.Error()}
		return res
	}

	return res
}
