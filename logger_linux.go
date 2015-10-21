// +build linux

package main

import (
	"fmt"
	"log/syslog"
)

type Logger struct {
	w *syslog.Writer
}

func NewLogger() (*Logger, error) {
	l := &Logger{}
	w, err := syslog.New(syslog.LOG_NOTICE, "qemu-ga")
	if err != nil {
		return nil, err
	}
	l.w = w
	return l, nil
}

func (l *Logger) Close() error {
	return l.w.Close()
}

func (l *Logger) Alert(msg string) error {
	return l.w.Alert(msg)
}

func (l *Logger) Alertf(f string, msg string) error {
	return l.Alert(fmt.Sprintf(f, msg))
}

func (l *Logger) Crit(msg string) error {
	return l.w.Crit(msg)
}

func (l *Logger) Critf(f string, msg string) error {
	return l.Crit(fmt.Sprintf(f, msg))
}

func (l *Logger) Debug(msg string) error {
	if !options.Verbose {
		return nil
	}
	return l.w.Debug(msg)
}

func (l *Logger) Debugf(f string, msg string) error {
	return l.Debug(fmt.Sprintf(f, msg))
}

func (l *Logger) Emerg(msg string) error {
	return l.w.Emerg(msg)
}

func (l *Logger) Emergf(f string, msg string) error {
	return l.Emerg(fmt.Sprintf(f, msg))
}

func (l *Logger) Error(msg string) error {
	return l.w.Err(msg)
}

func (l *Logger) Errorf(f string, msg string) error {
	return l.Error(fmt.Sprintf(f, msg))
}

func (l *Logger) Info(msg string) error {
	return l.w.Info(msg)
}

func (l *Logger) Infof(f string, msg string) error {
	return l.Info(fmt.Sprintf(f, msg))
}

func (l *Logger) Notice(msg string) error {
	return l.w.Notice(msg)
}

func (l *Logger) Noticef(f string, msg string) error {
	return l.Notice(fmt.Sprintf(f, msg))
}

func (l *Logger) Warn(msg string) error {
	return l.w.Warning(msg)
}

func (l *Logger) Warnf(f string, msg string) error {
	return l.Warn(fmt.Sprintf(f, msg))
}
