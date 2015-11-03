/*

guest-exec command have two versions.

First:
	{ "execute": "guest-exec", "arguments": {
		"command": string
		}
	}

Second:
	{ "execute": "guest-exec", "arguments": {
		"path": string,
		"arg": string,
		"env": string,
		"input": string,
		"capture-output": bool
		}
	}

*/
package main // import "github.com/vtolstov/qemu-ga"
