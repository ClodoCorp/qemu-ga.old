package main

import (
	"bufio"
	"errors"
	"io"
	"runtime"
	"time"
)

const ReaderBufferSize = 4 * 1024

var ReaderErrTimeout = errors.New("timeout")

type TimeoutReader struct {
	b  *bufio.Reader
	t  time.Duration
	ch <-chan error
}

func NewTimeoutReader(r io.Reader) *TimeoutReader {
	return &TimeoutReader{b: bufio.NewReaderSize(r, ReaderBufferSize), t: -1}
}

// SetTimeout sets the timeout for all future Read calls as follows:
//
// t < 0  -- block
// t == 0 -- poll
// t > 0  -- timeout after t
func (r *TimeoutReader) SetTimeout(t time.Duration) time.Duration {
	prev := r.t
	r.t = t
	return prev
}

func (r *TimeoutReader) Read(b []byte) (n int, err error) {
	if r.ch == nil {
		if r.t < 0 || r.b.Buffered() > 0 {
			return r.b.Read(b)
		}
		ch := make(chan error, 1)
		r.ch = ch
		go func() {
			_, err := r.b.Peek(1)
			ch <- err
		}()
		runtime.Gosched()
	}
	if r.t < 0 {
		err = <-r.ch // Block
	} else {
		select {
		case err = <-r.ch: // Poll
		default:
			if r.t == 0 {
				return 0, ReaderErrTimeout
			}
			select {
			case err = <-r.ch: // Timeout
			case <-time.After(r.t):
				return 0, ReaderErrTimeout
			}
		}
	}
	r.ch = nil
	if r.b.Buffered() > 0 {
		n, _ = r.b.Read(b)
	}
	return
}
