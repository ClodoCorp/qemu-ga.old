package main

var l *Logger

type FileSystem struct {
	Device  string
	Path    string
	Type    string
	Options []string
}
