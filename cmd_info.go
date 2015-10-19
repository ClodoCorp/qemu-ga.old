package main

import "encoding/json"

var cmdInfo = &Command{
	Name: "guest-info",
	Func: fnInfo,
}

var (
	Version   string
	BuildTime string
)

func init() {
	commands = append(commands, cmdInfo)
}

func fnInfo(m json.RawMessage) json.RawMessage {
	type command struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
		Success bool   `json:"success-response"`
	}

	res := struct {
		Return struct {
			Version  string    `json:"version"`
			Commands []command `json:"supported_commands"`
		} `json:"return"`
	}{}
	res.Return.Version = Version

	for _, cmd := range commands {
		res.Return.Commands = append(res.Return.Commands, command{Name: cmd.Name, Enabled: true, Success: true})
	}
	buf, err := json.Marshal(res)
	if err != nil {

	}
	return json.RawMessage(buf)
}
