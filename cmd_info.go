package main

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

func fnInfo(req *Request) *Response {
	res := &Response{}

	type command struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
		Success bool   `json:"success-response"`
	}

	info := struct {
		Version  string    `json:"version"`
		Commands []command `json:"supported_commands"`
	}{Version: Version}

	for _, cmd := range commands {
		info.Commands = append(info.Commands, command{Name: cmd.Name, Enabled: true, Success: true})
	}
	res.Return = info
	res.Id = req.Id
	return res
}
