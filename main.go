package main

//go:generate go run generate.go

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/vtolstov/qemu-ga/qga"
)

var (
	l *qga.Logger
)

func main() {
	var err error

	parser := flags.NewParser(&options, flags.PrintErrors)
	if _, err := parser.Parse(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	if options.Help {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	if options.Version {
		fmt.Printf("%s\n", qga.Version)
		os.Exit(0)
	}

	if options.Fork {
		if err = fork(); err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	/*
		l, err = qga.NewLogger()
		if err != nil {
			fmt.Printf(err.Error())
		}
	*/

	if err = master(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
