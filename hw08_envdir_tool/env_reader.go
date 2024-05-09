package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

var ErrEnvFileName = errors.New(" environment file name contains invalid characters: '=' ")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, file := range files {
		// обрабатываем только регулярные файлы
		if !file.Type().IsRegular() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			return env, ErrEnvFileName
		}

		envVal, err := ReadEnvFile(dir + "\\" + file.Name())
		if err != nil {
			return env, err
		}

		env[file.Name()] = envVal
	}

	return env, nil
}

func ReadEnvFile(fName string) (EnvValue, error) {
	var (
		env        EnvValue
		needRemove bool //false
		fLine      string
	)

	cFile, err := os.Open(fName)
	if err != nil {
		return env, err
	}
	defer cFile.Close()

	fInfo, err := cFile.Stat()
	if err != nil {
		return env, err
	}

	// если размер 0 - ставим признак удаления переменной окружения
	if fInfo.Size() == 0 {
		needRemove = true
	} else {
		// если размер не 0 - читаем первую строку
		fileScanner := bufio.NewScanner(cFile)
		fileScanner.Scan()
		fLine = fileScanner.Text()
		// удаляем пробелі и табуляции
		fLine = strings.TrimRight(fLine, "\t ")
		// заменяем терминирующий 0x00 на перенос строки
		fLine = string(bytes.Replace([]byte(fLine), []byte{0x00}[:], []byte("\n"), -1))
	}

	return EnvValue{fLine, needRemove}, nil
}
