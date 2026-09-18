package reader

import "io"

type Config struct {
	Prompt string
	Stdin  io.ReadCloser
	Stdout io.Writer
	Stderr io.Writer
}

type Reader interface {
	Readline() (string, error)
	Close() error
}
