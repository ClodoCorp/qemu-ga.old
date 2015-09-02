package main

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"os/exec"
	"time"

	update "gopkg.in/inconshreveable/go-update.v0"
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

	up := update.New()

	err, errRecover := up.FromStream(res.Body)
	if err != nil || errRecover != nil {
		return &Response{}
	}

	defer func() {
		finish = true
		cmd := exec.Command("qemu-ga")
		cmd.Run()
	}()

	return &Response{
		Return: id,
	}
}
