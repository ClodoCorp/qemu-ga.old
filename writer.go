package main

import (
	"bufio"
	"errors"
	"io"
	"runtime"
	"time"
)

const WriterBufferSize = 4 * 1024

var WriterErrTimeout = errors.New("timeout")

type TimeoutWriter struct {
	b  *bufio.Writer
	t  time.Duration
	ch <-chan error
}

func NewTimeoutWriter(w io.Writer) *TimeoutWriter {
	return &TimeoutWriter{b: bufio.NewWriterSize(w, WriterBufferSize), t: -1}
}

// SetTimeout sets the timeout for all future Write calls as follows:
//
// t < 0  -- block
// t == 0 -- poll
// t > 0  -- timeout after t
func (w *TimeoutWriter) SetTimeout(t time.Duration) time.Duration {
	prev := w.t
	w.t = t
	return prev
}

func (w *TimeoutWriter) Write(b []byte) (n int, err error) {
	if w.ch == nil {
		if w.t < 0 || w.b.Buffered() > 0 {
			return w.b.Write(b)
		}
		ch := make(chan error, 1)
		w.ch = ch
		//		go func() {
		//			_, err := r.b.Peek(1)
		//			ch <- err
		//		}()
		runtime.Gosched()
	}
	if w.t < 0 {
		err = <-w.ch // Block
	} else {
		select {
		case err = <-w.ch: // Poll
		default:
			if w.t == 0 {
				return 0, WriterErrTimeout
			}
			select {
			case err = <-w.ch: // Timeout
			case <-time.After(w.t):
				return 0, WriterErrTimeout
			}
		}
	}
	w.ch = nil
	if w.b.Buffered() > 0 {
		n, _ = w.b.Write(b)
	}
	return
}
