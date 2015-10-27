package main

var cmdInfo = &Command{
	Name:    "guest-info",
	Func:    fnInfo,
	Enabled: true,
	Returns: true,
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

	info := struct {
		Version  string     `json:"version"`
		Commands []*Command `json:"supported_commands"`
	}{Version: Version, Commands: commands}

	res.Return = info
	res.Id = req.Id
	return res
}
