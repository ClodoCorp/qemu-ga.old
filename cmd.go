package main

// Command struct contains supported commands
type Command struct {
	Name string
	Func func(map[string]interface{}) interface{}
}

var commands = []*Command{}
