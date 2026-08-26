package shell

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func (s *Shell) Run() int {
	for {
		if _, err := fmt.Fprint(s.out, "$ "); err != nil {
			return 1
		}

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

		if result := s.execute(command); result.Exit {
			return result.Status
		}
	}
}
