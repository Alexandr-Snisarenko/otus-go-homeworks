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
	Connect(context.Context) error
	Close() error
	IsActive() bool
	Done() <-chan struct{}
	// Send() error
	// Receive() error
}

type telnetClient struct {
	mu        sync.Mutex
	active    bool
	done      chan struct{}
	address   string
	timeout   time.Duration
	inReader  io.Reader
	outWriter io.Writer
	errWriter io.Writer
	conn      net.Conn
}

func (t *telnetClient) Connect(ctx context.Context) error {
	if t.active {
		return ErrConnectionAlreadyActive
	}

	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	t.active = true
	t.done = make(chan struct{})

	go func() {
		defer t.Close()

		t.inBuf2outBuff(t.inReader, t.conn)
		t.inBuf2outBuff(t.conn, t.outWriter)

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
	return nil
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

func (t *telnetClient) Done() <-chan struct{} {
	return t.done
}

func NewTelnetClient(address string, timeout time.Duration, inReader io.Reader, outWriter io.Writer, errWriter io.Writer) TelnetClient {
	return &telnetClient{active: false, address: address, timeout: timeout, inReader: inReader, outWriter: outWriter, errWriter: errWriter}
}

func (t *telnetClient) inBuf2outBuff(inBuf io.Reader, outBuf io.Writer) {
	go func() {
		reader := bufio.NewReader(inBuf)
		for {
			select {
			case <-t.done:
				return
			default:
				text, err := reader.ReadString('\n')
				if err != nil {
					t.errWriter.Write([]byte(err.Error() + "\n"))
					t.Close()
					return
				}

				_, err = outBuf.Write([]byte(text))
				if err != nil {
					t.errWriter.Write([]byte(err.Error() + "\n"))
					t.Close()
					return
				}
			}
		}
	}()
}
