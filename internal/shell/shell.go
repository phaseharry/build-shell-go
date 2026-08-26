package shell

import (
	"bufio"
	"io"
)

type builtin func(args []string) Result

// Result describes the outcome of running one command.
type Result struct {
	Status int
	Exit   bool
}

// Shell owns the state and dependencies of a shell session.
type Shell struct {
	in         io.Reader
	out        io.Writer
	errOut     io.Writer
	reader     *bufio.Reader
	builtins   map[string]builtin
	lastStatus int
}

// New constructs a shell using the supplied standard streams.
func New(
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
) *Shell {
	sh := &Shell{
		in:     in,
		out:    out,
		errOut: errOut,
		reader: bufio.NewReader(in),
	}
	sh.builtins = sh.builtinRegistry()

	return sh
}
