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
		return &Response{}
	}
	defer res.Body.Close()
	err = update.Apply(res.Body, &update.Options{TargetMode: os.FileMode(0700)})
	if err != nil {
		return &Response{}
	}

	defer func() {
		cmd := exec.Command("qemu-ga")
		cmd.Env = append(cmd.Env, fmt.Sprintf("PARENT=%d", os.Getpid()))
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Noctty: true, Setpgid: true}
		cmd.Start()
	}()

	return &Response{
		Return: id,
	}
}
