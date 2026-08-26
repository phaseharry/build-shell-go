package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *Shell) builtRegistry() map[string]builtin {
	return map[string]builtin{
		"exit": s.exit,
		"echo": s.echo,
		"type": s.typeCommand,
		"pwd":  s.pwd,
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

// using typeCommand instead of type since type is a keyword
func (s *Shell) typeCommand(args []string) Result {
	for _, command := range args {
		_, ok := s.builtins[command]
		if ok {
			fmt.Fprintf(s.out, "%s is a shell builtin\n", command)
			continue
		}
		// TODO, implement PATH lookup instead of using exec.LookPath
		path, err := exec.LookPath(command)
		if err != nil {
			fmt.Fprintf(s.out, "%s: not found\n", command)
		} else {
			fmt.Fprintf(s.out, "%s is %s\n", command, path)
		}
	}

	return Result{
		Exit:   false,
		Status: 0,
	}
}

func (s *Shell) pwd(_args []string) Result {
	directory, err := os.Getwd()
	if err != nil {
		return Result{
			Exit:   false,
			Status: 1,
		}
	}
	fmt.Fprintln(s.out, directory)

	return Result{
		Exit: false,
		Status: 0,
	}
}
