package main

import (
	"fmt"
	"os"
)

func main() {

	path := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = RunCmd(cmd, env); err != nil {
		fmt.Println(err)
	}
}
