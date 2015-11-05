package qga

import (
	"encoding/json"
	"time"
)

// Channel interface provide communication channel with qemu-ga
type Channel interface {
	DialTimeout(string, time.Duration) error
	Close() error
	Poll() error
}

// Request struct used to parse incoming request
type Request struct {
	Execute string          `json:"execute"`
	RawArgs json.RawMessage `json:"arguments,omitempty"`
	Id      string          `json:"id,omitempty"`
}

// Error struct used to indicate error when processing command
type Error struct {
	Class  string `json:"class,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Bufb64 string `json:"bufb64,omitempty"`
	Code   int    `json:"code,omitempty"`
}

// Response struct used to encode response from command
type Response struct {
	Return interface{} `json:"return,omitempty"`
	Error  *Error      `json:"error,omitempty"`
	Id     string      `json:"id,omitempty"`
}
