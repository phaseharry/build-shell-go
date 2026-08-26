package shell

import (
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func (s *Shell) execute(command parser.Command) Result {
	builtinFnc, ok := s.builtins[command.Name]
	if ok {
		return builtinFnc(command.Args)
	}

	_, err := exec.LookPath(command.Name)
	if err != nil {
		fmt.Fprintf(s.out, "%s: command not found\n", command.Name)
		return Result{}
	}

	cmd := exec.Command(command.Name, command.Args...)
	cmd.Stdin = s.in
	cmd.Stdout = s.out
	cmd.Stderr = s.errOut

	if err := cmd.Run(); err != nil {
		return Result{
			Exit:   false,
			Status: 1,
		}
	}

	return Result{}
}
