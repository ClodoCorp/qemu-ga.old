package main

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
	"syscall"
	"time"
)

var cmdUpdate = &Command{
	Name: "guest-agent-update",
	Func: fnUpdate,
}

func init() {
	commands = append(commands, cmdUpdate)
}

func fnUpdate(req *Request) *Response {
	res := &Response{}
	var r io.ReadCloser

	httpTransport := &http.Transport{
		Dial:            (&net.Dialer{DualStack: true}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: httpTransport, Timeout: 20 * time.Second}

	update := struct {
		Path string `json:"path"`
	}{}

	err := json.Unmarshal(req.RawArgs, &update)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}

	u, err := url.Parse(update.Path)
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	switch u.Scheme {
	case "http", "https":
		hres, err := httpClient.Get(update.Path)
		if err != nil {
			res.Error = &Error{Code: -1, Desc: err.Error()}
			return res
		}
		r = hres.Body
	case "file":
		r, err = os.Open(u.Path)
		if err != nil {
			res.Error = &Error{Code: -1, Desc: err.Error()}
			return res
		}
	default:
		res.Error = &Error{Code: -1, Desc: fmt.Sprintf("invalid path %s", u)}
		return res
	}
	defer r.Close()

	dirname, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	filename := fmt.Sprintf(".%s", filepath.Base(os.Args[0]))
	w, err := os.OpenFile(filepath.Join(dirname, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0755))
	if err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	_, err = io.Copy(w, r)
	if err != nil {
		defer w.Close()
		defer os.Remove(filepath.Join(dirname, filename))
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	w.Sync()
	w.Close()

	if err = os.Rename(filepath.Join(dirname, filename), filepath.Join(dirname, filepath.Base(os.Args[0]))); err != nil {
		res.Error = &Error{Code: -1, Desc: err.Error()}
		return res
	}
	time.Sleep(2 * time.Second)
	defer func() {
		cmd := exec.Command(os.Args[0])
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Noctty: false, Setpgid: false, Foreground: false}

		err = cmd.Start()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()

	ret := struct {
		Id int `json:"-"`
	}{}
	res.Return = ret
	return res
}
