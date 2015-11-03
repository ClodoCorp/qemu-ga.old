/*

guest-exec command have two versions.

First:
	{ "execute": "guest-exec", "arguments": {
		"command": string // required, base64 encoded command name to execute with args including newline
		}
	}

Second:
	{ "execute": "guest-exec", "arguments": {
		"path": string, // required, command name to execute
		"arg": string, // optional, arguments to executed command
		"env": string, // optional, environment to executed command
		"input": string, // optional, base64 encoded string
		"capture-output": bool // optional, capture stdout/stderr
		}
	}

*/
package main // import "github.com/vtolstov/qemu-ga"
