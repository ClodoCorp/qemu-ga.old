package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	update "github.com/inconshreveable/go-update"
)

var cmdUpdate = &Command{
	Name: "guest-agent-update",
	Func: fnUpdate,
}

func init() {
	commands = append(commands, cmdUpdate)
}

func fnUpdate(d map[string]interface{}) interface{} {

	httpTransport := &http.Transport{
		Dial:            (&net.Dialer{DualStack: true}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: httpTransport, Timeout: 20 * time.Second}

	id, _ := (d["id"].(json.Number)).Int64()
	path := d["path"].(string)

	res, err := httpClient.Get(path)
	if err != nil {
		return &Response{Return: err.Error()}
	}
	defer res.Body.Close()
	err = update.Apply(res.Body, &update.Options{TargetMode: os.FileMode(0755)})
	if err != nil {
		return &Response{Return: err.Error()}
	}

	defer func() {
		cmd := exec.Command("qemu-ga")
		cmd.Env = append(cmd.Env, fmt.Sprintf("PARENT=%d", ppid))
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Noctty: false, Setpgid: false, Foreground: false}

		err = cmd.Start()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()

	return &Response{
		Return: id,
	}
}
