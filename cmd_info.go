package main

var cmdInfo = &Command{
	Name: "guest-info",
	Func: fnInfo,
}

func init() {
	commands = append(commands, cmdInfo)
}

func fnInfo(d map[string]interface{}) interface{} {
	type command struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	}

	type response struct {
		Version  string    `json:"version"`
		Commands []command `json:"supported_commands"`
	}
	res := &response{Version: "1.5.2"}

	for _, cmd := range commands {
		res.Commands = append(res.Commands, command{Name: cmd.Name, Enabled: true})
	}
	return &Response{Return: res}
}
