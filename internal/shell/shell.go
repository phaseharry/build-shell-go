package shell

import (
	"bufio"
	"io"
)

type builtin func(args []string) Result

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
