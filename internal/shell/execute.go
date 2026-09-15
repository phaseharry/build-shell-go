package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func (s *Shell) execute(command parser.Command) Result {
	std := Streams{In: s.in, Out: s.out, Err: s.errOut}

	for _, redirect := range command.Redirects {
		// open file flags to Write ONLY and Create if not exist initial
		flags := os.O_WRONLY | os.O_CREATE
		// add Append flag if the redirect is a >>, else it's an overwrite flag
		if redirect.Append {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		file, err := os.OpenFile(redirect.Target, flags, 0644)
		if err != nil {
			fmt.Fprintf(s.errOut, "%s: No such file or directory\n", redirect.Target)
			return Result{Status: 1}
		}
		defer file.Close()

		switch redirect.FileDescriptor {
		case 1:
			std.Out = file
		case 2:
			std.Err = file
		default:
			std.Out = file
		}

	}

	builtinFnc, ok := s.builtins[command.Name]
	if ok {
		return builtinFnc(command.Args, std)
	}

	_, err := exec.LookPath(command.Name)
	if err != nil {
		fmt.Fprintf(s.out, "%s: command not found\n", command.Name)
		return Result{}
	}

	cmd := exec.Command(command.Name, command.Args...)
	cmd.Stdin = std.In
	cmd.Stdout = std.Out
	cmd.Stderr = std.Err

	if err := cmd.Run(); err != nil {
		return Result{
			Exit:   false,
			Status: 1,
		}
	}

	return Result{}
}
