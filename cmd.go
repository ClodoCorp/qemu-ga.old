package main

import "encoding/json"

// Command struct contains supported commands
type Command struct {
	Name string
	Func func(json.RawMessage) json.RawMessage
}

var commands = []*Command{}
