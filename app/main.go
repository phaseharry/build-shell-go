package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/shell"
)

func main() {
	sh := shell.New(os.Stdin, os.Stdout, os.Stderr)
	sh.Run()
}
