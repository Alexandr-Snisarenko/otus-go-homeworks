package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

func in2outRoutine(ctx context.Context, chIn <-chan string, wg *sync.WaitGroup, iOut io.Writer) error {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return nil
		case text, ok := <-chIn:
			if !ok {
				return nil
			}
			_, err := iOut.Write([]byte(text))
			if err != nil {
				return fmt.Errorf("error write to Output: %w ", err)
			}
		}
	}
}

func read2chan(iRead io.Reader) <-chan string {
	out := make(chan string)
	reader := bufio.NewReader(iRead)
	go func() {
		defer close(out)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("..EOF")
					//					return
				}
				break
			}
			out <- text
		}
	}()
	return out
}

func main() {
	var timeout time.Duration

	// логер. debug mode
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(log)

	// таймаут подключения к серверу. по умолчанию 10 сек.
	pflag.DurationVar(&timeout, "timeout", 10, "timeout in second waiting for connection establishment")
	pflag.Parse()

	// аргументы командной строки (адрес и порт) должно быть 2. если не 2 - ошибка.
	args := pflag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Arguments count is wrong. Expected 2, recieved %d", len(args))
		return
	}

	// подключение к серверу
	addr := args[0] + ":" + args[1]
	conn, err := net.DialTimeout("tcp", addr, timeout*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to host %s. Error: %s", addr, err.Error())
		return
	} else {
		defer conn.Close()
	}

	// контекст с отменой.
	ctx, cancel := context.WithCancel(context.Background())

	chStdIn := read2chan(os.Stdin)
	chConn := read2chan(conn)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// горутиа чтения данных с сервера. читаем в stdout
	go func() {
		err := in2outRoutine(ctx, chConn, wg, os.Stdout)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		cancel()
	}()

	wg.Add(1)
	// горутина записи данных в сервер. пишем из stdin
	go func() {
		err := in2outRoutine(ctx, chStdIn, wg, conn)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		cancel()
	}()

	wg.Wait()
}
