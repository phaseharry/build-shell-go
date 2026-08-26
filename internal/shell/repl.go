package shell

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func (s *Shell) Run() int {
	for {
		fmt.Fprint(s.out, "$ ")

		line, err := s.reader.ReadString('\n')
		if err != nil {
			fmt.Fprint(s.out, err.Error())
			return 1
		}

		command := parser.ParseLine(line)

		// handle empty line by just keeping the repl going
		if command.Name == "" {
			continue
		}

		fnc, ok := s.builtins[command.Name]
		if !ok {
			fmt.Fprintf(s.out, "%s: command not found\n", command.Name)
			continue
		}

		result := fnc([]string{command.Args})
		if result.Exit {
			return result.Status
		}
	}
}
