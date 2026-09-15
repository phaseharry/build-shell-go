package shell

import (
	"bufio"
	"io"
)

type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type builtin func(args []string, std Streams) Result

type Result struct {
	Status int
	Exit   bool
}

type Shell struct {
	in       io.Reader
	out      io.Writer
	errOut   io.Writer
	reader   *bufio.Reader
	builtins map[string]builtin
}

func New(in io.Reader, out io.Writer, errOut io.Writer) *Shell {
	sh := &Shell{
		in:     in,
		out:    out,
		errOut: errOut,
		reader: bufio.NewReader(in),
	}
	sh.builtins = sh.builtRegistry()
	return sh
}
