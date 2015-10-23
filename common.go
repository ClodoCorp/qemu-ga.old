package main

import "os"

var l *Logger

type FileSystem struct {
	Device  string
	Path    string
	Type    string
	Options []string
}

var openFiles map[int]*os.File

func init() {
	openFiles = make(map[int]*os.File)
}
