package main

import (
	"errors"
	"os"
	"strings"
	"syscall"
)

var ErrOsEnvironmentVariableFormat = errors.New(" environment variable string don't contain '=' ")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) error {
	var (
		newEnv []string
		err    error
	)
	if newEnv, err = parseOsEnv(env); err != nil {
		return err
	}

	err = syscall.Exec(cmd[0], cmd[1:], newEnv)
	if err != nil {
		return err
	}

	return nil
}

// функция парсировки текущего системного списка переменных окружения
// и формирования нового списка переменных окружения
// с учетом входного словаря переменных "на замену"
// логика следующая:
// 1 - переносим в новый список все переменные из системного списка, которых нет во входном словаре.
// 2 - обрабатываем входной словарь и добавляем в новый список те переменные, у которых NeedRemove == false.
func parseOsEnv(env Environment) ([]string, error) {
	// получаем список текущих переменных окружения
	osEnv := os.Environ()
	// массив для новых переменных окружения (примерно такой же длинны)
	newEnv := make([]string, 0, len(osEnv))

	// проходим по текущему системному списку переменных окружения
	for _, osEnv := range osEnv {
		// определяем имя переменной по символу "="
		// если "=" нет - unexpected error (такого быть не может по идее)
		before, _, ok := strings.Cut(osEnv, "=")
		if !ok {
			return nil, ErrOsEnvironmentVariableFormat
		}

		// ищем переменную окружения в нашем словаре
		// если нашли - пропускаем (они обработаются отдельным блоком)
		// если не нашли - переносим в новый массив as is
		if _, ok := env[before]; !ok {
			newEnv = append(newEnv, osEnv)
		}
	}

	// проходим по нашему словарю переменных окружения
	// если признак NeedRemove == false - добавляем переменную в новый список
	for eVal := range env {
		if !env[eVal].NeedRemove {
			newEnv = append(newEnv, eVal+"="+env[eVal].Value)
		}
	}

	return newEnv, nil
}
