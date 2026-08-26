package shell

import (
	"fmt"
	"strings"
)

func (s *Shell) Run() int {
	for {
		fmt.Fprint(s.out, "$ ")

		line, err := s.reader.ReadString('\n')
		if err != nil {
			fmt.Fprint(s.out, err.Error())
			return 1
		}

		// using string.Fields to split all strings by whitespace, removing them all

		tokens := strings.Fields(line)
		// handle an empty line by just keeping the repl going
		if len(tokens) == 0 {
			continue
		}
		command, args := tokens[0], tokens[1:]

		fnc, ok := s.builtins[command]
		if !ok {
			fmt.Fprintf(s.out, "%s: command not found\n", command)
			continue
		}

		result := fnc(args)
		if result.Exit {
			return result.Status
		}
	}
}
