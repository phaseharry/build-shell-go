package shell

import (
	"io"

	"github.com/codecrafters-io/shell-starter-go/internal/reader"
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
	reader   reader.Reader
	builtins map[string]builtin
}

func New(in io.ReadCloser, out io.Writer, errOut io.Writer) *Shell {
	sh := &Shell{
		in:     in,
		out:    out,
		errOut: errOut,
		reader: reader.New(
			reader.Config{
				Prompt: "$ ",
				Stdin:  in,
				Stdout: out,
				Stderr: errOut,
			}),
	}
	sh.builtins = sh.builtRegistry()
	return sh
}
