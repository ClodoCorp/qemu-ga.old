package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
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
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	l, err = NewLogger()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	if err = master(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
