package main

import (
	"encoding/json"
	"time"
)

type Channel interface {
	DialTimeout(string, time.Duration) error
	Close() error
	Poll() error
}

type Request struct {
	Execute string          `json:"execute"`
	RawArgs json.RawMessage `json:"arguments,omitempty"`
	Id      string          `json:"id,omitempty"`
}

type Error struct {
	Class  string `json:"class,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Bufb64 string `json:"bufb64,omitempty"`
	Code   int    `json:"code,omitempty"`
}

type Response struct {
	Return interface{} `json:"return,omitempty"`
	Error  *Error      `json:"error,omitempty"`
	Id     string      `json:"id,omitempty"`
}
