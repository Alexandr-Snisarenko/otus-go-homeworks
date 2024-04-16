package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileExists            = errors.New("destination file already exists")
)

func Copy(fromPath, toPath string, offset, limit int64, rewrite bool) error {
	var (
		inFile, outFile         *os.File
		inFileInfo, outFileInfo fs.FileInfo
		err                     error
	)

	/// проверки параметров файлов ///
	inFileInfo, err = os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("can't open file %s : %w", fromPath, err)
	}
	if mode := inFileInfo.Mode(); !mode.IsRegular() {
		return ErrUnsupportedFile
	}

	// если свдиг больше файла - ошибка
	if offset >= inFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// если целевой файл существует - проверяем можно ли его перезаписать (флаг rewrite и тип файла)
	outFileInfo, err = os.Stat(toPath)
	if err == nil {
		// проверяем флаг разрешения перезаписи
		if !rewrite {
			return ErrFileExists
		}
		// проверяем тип файла обрабатываем только регулярные файлы
		if mode := outFileInfo.Mode(); !mode.IsRegular() {
			return fmt.Errorf("can't rewrite destination file %s: %w", toPath, ErrUnsupportedFile)
		}
	}

	/// обработка файлов ///
	inFile, err = os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("can't open source file: %w", err)
	}
	defer inFile.Close()

	outFile, err = os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can't create destination file: %w", err)
	}
	defer outFile.Close()

	fmt.Println("Copying ", inFile.Name(), " to ", outFile.Name())

	// если задан сдвиг - сдвигаем ридер
	if offset > 0 {
		inFile.Seek(offset, io.SeekStart)
	}

	// если лимит не задан, или лимит, с учетом сдвига, больше размера файла, то limit (размер копируемых данных) = размер файла - сдвиг
	if limit == 0 || offset+limit > inFileInfo.Size() {
		limit = inFileInfo.Size() - offset
	}

	bar := progressbar.DefaultBytes(limit)

	buf := make([]byte, 1024) //  буфер обмена
	exit := false             // флаг завершения копирования
	cntAll := int64(0)        // счетчик скопированных байт
	for !exit {
		cnt, err := inFile.Read(buf)
		cntAll += int64(cnt)

		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error read from source file: %w", err)
			}
			exit = true
		}

		// проверяем на limit. если лимит првышен, обрезаем крайнюю порцию по размеру превышения и ставим флаг "на выход"
		if cntAll >= limit {
			cnt -= int(cntAll - limit)
			exit = true
		}

		_, err = outFile.Write(buf[:cnt])
		if err != nil {
			return fmt.Errorf("error write to destination file: %w", err)
		}
		// актуализируем прогрессбар
		bar.Add(cnt)
	}

	fmt.Println("Copying completed")
	return nil

}
