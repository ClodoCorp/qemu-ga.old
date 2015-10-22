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
		Version  string     `json:"version"`
		Commands []*Command `json:"supported_commands"`
	}{Version: Version}

	info.Commands = commands

	res.Return = info
	res.Id = req.Id
	return res
}
