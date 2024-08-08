package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/pflag"
)

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
		fmt.Fprintf(os.Stderr, "Arguments count is wrong. Expected 2, recieved %d \n", len(args))
		return
	}

	// переводим в секунды
	timeout = timeout * time.Second
	// формируем строку адрес:порт из аргументов командной сроки
	addr := args[0] + ":" + args[1]

	// создаем объект телнет клиента. в качетсве потоков данных указываем std[in|out|err]
	t := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout, os.Stderr)

	// пробуем подключиться. если ошибка - выводим её в stderr
	err := t.Connect(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to host %s. Error: %s\n", addr, err.Error())
		return
	}

	// ждем завершения работы телнет клиента
	<-t.Done()

}
