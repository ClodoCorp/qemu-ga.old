package main

import (
	"fmt"
	"time"

	"github.com/vtolstov/qemu-ga/qga"
)

func slave() error {
	var ch qga.Channel
	var err error

	switch options.Method {
	case "virtio-serial":
		if ch, err = NewVirtioChannel(); err != nil {
			return err
		}
		err = ch.DialTimeout(options.Path, time.Minute)
		/*
			case "isa-serial":
				if ch, err = NewIsaChannel(); err != nil {
					return err
				}
				err = ch.DialTimeout(options.Path, time.Minute)
		*/
	default:
		return fmt.Errorf("unsupported method %s", options.Method)
	}
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Poll()
}
