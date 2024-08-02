package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

var (
	ErrConnectionAlreadyActive = errors.New("connection is already active")
)

type TelnetClient interface {
	Connect(context.Context) (<-chan struct{}, error)
	Close() error
	IsActive() bool
	// Send() error
	// Receive() error
}

type telnetClient struct {
	mu        sync.Mutex
	active    bool
	done      chan struct{}
	address   string
	timeout   time.Duration
	inReader  io.ReadCloser
	outWriter io.Writer
	errWriter io.Writer
	conn      net.Conn
}

func (t *telnetClient) Connect(ctx context.Context) (<-chan struct{}, error) {
	if t.active {
		return nil, ErrConnectionAlreadyActive
	}

	conn, err := net.DialTimeout("tcp", t.address, t.timeout*time.Second)
	if err != nil {
		return nil, err
	}
	t.conn = conn
	t.active = true
	t.done = make(chan struct{})

	go func() {
		defer t.Close()
		chIn := reader2chan(t.done, t.inReader, t.errWriter, t.Close)
		chConn := reader2chan(t.done, t.conn, t.errWriter, t.Close)

		chan2writer(t.done, chIn, t.conn, t.errWriter, t.Close)
		chan2writer(t.done, chConn, t.outWriter, t.errWriter, t.Close)

		for {
			select {
			case <-t.done:
				return
			case <-ctx.Done():
				t.Close()
				return
			}
		}

	}()

	return t.done, nil
}

func (t *telnetClient) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.active {
		return nil
	}
	t.active = false
	close(t.done)

	return t.conn.Close()
}

func (t *telnetClient) IsActive() bool {
	return t.active
}

func NewTelnetClient(address string, timeout time.Duration, inReader io.ReadCloser, outWriter io.Writer, errWriter io.Writer) TelnetClient {
	return &telnetClient{active: false, address: address, timeout: timeout, inReader: inReader, outWriter: outWriter, errWriter: errWriter}
}

func chan2writer(done <-chan struct{}, chIn <-chan string, bufOut io.Writer, bufErr io.Writer, closeConn func() error) {
	go func() {
		for {
			select {
			case <-done:
				return
			case text, ok := <-chIn:
				if !ok {
					return
				}
				_, err := bufOut.Write([]byte(text))
				if err != nil {
					bufErr.Write([]byte(err.Error() + "\n"))
					closeConn()
					return
				}
			}
		}
	}()
}

func reader2chan(done <-chan struct{}, bufRead io.Reader, bufErr io.Writer, closeConn func() error) <-chan string {
	out := make(chan string)
	reader := bufio.NewReader(bufRead)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			default:
				text, err := reader.ReadString('\n')
				if err != nil {
					bufErr.Write([]byte(err.Error() + "\n"))
					closeConn()
					return
				}
				out <- text
			}
		}
	}()
	return out
}
