package qga

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Logger struct {
	w *http.Client
}

func NewLogger(c *http.Client) (*Logger, error) {
	l := &Logger{}
	if c == nil {
		httpTransport := &http.Transport{
			Dial:            (&net.Dialer{DualStack: true}).Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		dt, err := time.ParseDuration("10s")
		if err != nil {
			return nil, err
		}
		l.w = &http.Client{Transport: httpTransport, Timeout: dt}
	} else {
		l.w = c
	}
	return l, nil
}

func (l *Logger) Close() error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Alert(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Alertf(f string, msg string) error {
	return l.Alert(fmt.Sprintf(f, msg))
}

func (l *Logger) Crit(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Critf(f string, msg string) error {
	return l.Crit(fmt.Sprintf(f, msg))
}

func (l *Logger) Debug(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Debugf(f string, msg string) error {
	return l.Debug(fmt.Sprintf(f, msg))
}

func (l *Logger) Emerg(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Emergf(f string, msg string) error {
	return l.Emerg(fmt.Sprintf(f, msg))
}

func (l *Logger) Error(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Errorf(f string, msg string) error {
	return l.Error(fmt.Sprintf(f, msg))
}

func (l *Logger) Info(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Infof(f string, msg string) error {
	return l.Info(fmt.Sprintf(f, msg))
}

func (l *Logger) Notice(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Noticef(f string, msg string) error {
	return l.Notice(fmt.Sprintf(f, msg))
}

func (l *Logger) Warn(msg string) error {
	if l.w == nil {
		return nil
	}
	return nil
}

func (l *Logger) Warnf(f string, msg string) error {
	return l.Warn(fmt.Sprintf(f, msg))
}
