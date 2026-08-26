package shell

import (
	"fmt"
	"strings"
)

func (s *Shell) builtRegistry() map[string]builtin {
	return map[string]builtin{
		"exit": s.exit,
		"echo": s.echo,
	}
}

func (s *Shell) exit(args []string) Result {
	// exiting the shell with status code of 0 because there was no error
	return Result{
		Exit:   true,
		Status: 0,
	}
}

func (s *Shell) echo(args []string) Result {
	fmt.Fprintf(s.out, "%s\n", strings.Join(args, " "))
	return Result{
		Exit:   false,
		Status: 0,
	}
}
