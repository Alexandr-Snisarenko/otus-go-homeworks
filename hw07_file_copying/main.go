package main

import (
	"flag"
	"fmt"
)

var (
	from, to  string
	limit, offset     int64
	rewrite bool
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.BoolVar(&rewrite, "rewrite", false, "rewrite file if exists ")
}

func main() {
	flag.Parse()
	// базовые проверки значений параметров
	switch {
	case from == "":
		fmt.Println("Name of file to read can't be empty")
		return
	case to == "":
		fmt.Println("Name of file to write can't be empty")
		return
	case limit < 0:
		fmt.Println("Limit can't be negative")
		return
	case offset < 0:
		fmt.Println("Offset can't be negative")
		return
	}

	if err := Copy(from, to, offset, limit, rewrite); err != nil {
		fmt.Println("Error:  ", err)
	}
}
