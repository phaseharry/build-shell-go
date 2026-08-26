package shell

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func (s *Shell) execute(command parser.Command) Result {
	if handler, ok := s.builtins[command.Name]; ok {
		return handler(command.Args)
	}

	path, ok := lookupExecutable(command.Name)
	if !ok {
		fmt.Fprintf(s.errOut, "%s: command not found\n", command.Name)
		return Result{Status: 127}
	}

	process := exec.Command(path, command.Args...)
	process.Stdin = s.in
	process.Stdout = s.out
	process.Stderr = s.errOut

	if err := process.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			status := exitError.ExitCode()
			if status < 0 {
				status = 1
			}
			return Result{Status: status}
		}

		fmt.Fprintf(s.errOut, "%s: %v\n", command.Name, err)
		return Result{Status: 126}
	}

	return Result{}
}
