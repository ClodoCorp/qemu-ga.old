package main

import (
	"encoding/json"
	"log"
)

var cmdSync = &Command{
	Name: "guest-sync",
	Func: fnSync,
}

func init() {
	commands = append(commands, cmdSync)
}

func fnSync(m json.RawMessage) json.RawMessage {
	req := struct {
		Id int `json:"id"`
	}{}
	res := struct {
		Id int `json:"return"`
	}{}
	err := json.Unmarshal(m, &req)
	if err != nil {
		log.Printf("aRRR %s\n", err.Error())
		return nil
	}
	res.Id = req.Id
	buf, err := json.Marshal(res)
	if err != nil {
		log.Printf("zRRR %s\n", err.Error())
		return nil
	}
	return json.RawMessage(buf)
}
