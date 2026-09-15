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
		"cd":   s.cd,
	}
}

func (s *Shell) exit(args []string, std Streams) Result {
	// exiting the shell with status code of 0 because there was no error
	return Result{
		Exit:   true,
		Status: 0,
	}
}

func (s *Shell) echo(args []string, std Streams) Result {
	fmt.Fprintf(std.Out, "%s\n", strings.Join(args, " "))
	return Result{
		Exit:   false,
		Status: 0,
	}
}

// using typeCommand instead of type since type is a keyword
func (s *Shell) typeCommand(args []string, std Streams) Result {
	for _, command := range args {
		_, ok := s.builtins[command]
		if ok {
			fmt.Fprintf(std.Out, "%s is a shell builtin\n", command)
			continue
		}
		// TODO, implement PATH lookup instead of using exec.LookPath
		path, err := exec.LookPath(command)
		if err != nil {
			fmt.Fprintf(std.Err, "%s: not found\n", command)
		} else {
			fmt.Fprintf(std.Out, "%s is %s\n", command, path)
		}
	}

	return Result{
		Exit:   false,
		Status: 0,
	}
}

func (s *Shell) pwd(args []string, std Streams) Result {
	directory, err := os.Getwd()
	if err != nil {
		return Result{
			Exit:   false,
			Status: 1,
		}
	}
	fmt.Fprintln(std.Out, directory)

	return Result{
		Exit:   false,
		Status: 0,
	}
}

func (s *Shell) cd(args []string, std Streams) Result {
	// defaulting to HOME directory if there are no arguments
	target := os.Getenv("HOME")

	if len(args) >= 1 {
		target = args[0]
	}

	// if the arg target is "~", reset it HOME since that's a shorthand for HOME
	if target == "~" {
		target = os.Getenv("HOME")
	}

	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(std.Err, "cd: %s: No such file or directory\n", target)
		return Result{
			Exit:   false,
			Status: 1,
		}
	}

	return Result{}
}
