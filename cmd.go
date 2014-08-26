package main

// Command struct contains commands supported by xenmgmd
type Command struct {
	Name string
	Func func(map[string]interface{}) interface{}
}

var commands = []*Command{}
