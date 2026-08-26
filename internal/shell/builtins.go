package shell

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

func (s *Shell) builtinRegistry() map[string]builtin {
	return map[string]builtin{
		"cd":   s.cd,
		"echo": s.echo,
		"exit": s.exit,
		"pwd":  s.pwd,
		"type": s.commandType,
	}
}

func (s *Shell) exit(args []string) Result {
	if len(args) == 0 {
		return Result{Exit: true}
	}
	if len(args) > 1 {
		fmt.Fprintln(s.errOut, "exit: too many arguments")
		return Result{Status: 1}
	}

	status, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(s.errOut, "exit: %s: numeric argument required\n", args[0])
		return Result{Status: 2, Exit: true}
	}

	return Result{Status: status, Exit: true}
}

func (s *Shell) echo(args []string) Result {
	if _, err := fmt.Fprintln(s.out, strings.Join(args, " ")); err != nil {
		return Result{Status: 1}
	}

	return Result{}
}

func (s *Shell) pwd(_ []string) Result {
	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(s.errOut, "pwd: %v\n", err)
		return Result{Status: 1}
	}
	if _, err := fmt.Fprintln(s.out, workingDirectory); err != nil {
		return Result{Status: 1}
	}

	return Result{}
}

func (s *Shell) cd(args []string) Result {
	target := "~"
	if len(args) > 0 {
		target = args[0]
	}

	target, err := expandHomePath(target)
	if err != nil {
		fmt.Fprintf(s.errOut, "cd: %v\n", err)
		return Result{Status: 1}
	}
	if err := os.Chdir(target); err != nil {
		message := err.Error()
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			message = pathError.Err.Error()
		}
		if errors.Is(err, fs.ErrNotExist) {
			message = "No such file or directory"
		}
		fmt.Fprintf(s.errOut, "cd: %s: %s\n", target, message)
		return Result{Status: 1}
	}

	return Result{}
}

func (s *Shell) commandType(args []string) Result {
	if len(args) == 0 {
		fmt.Fprintln(s.errOut, "type: missing operand")
		return Result{Status: 1}
	}

	status := 0
	for _, command := range args {
		if _, ok := s.builtins[command]; ok {
			fmt.Fprintf(s.out, "%s is a shell builtin\n", command)
			continue
		}
		if path, ok := lookupExecutable(command); ok {
			fmt.Fprintf(s.out, "%s is %s\n", command, path)
			continue
		}

		fmt.Fprintf(s.out, "%s: not found\n", command)
		status = 1
	}

	return Result{Status: status}
}
