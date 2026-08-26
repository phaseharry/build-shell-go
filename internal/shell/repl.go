package shell

import (
	"fmt"
	"strings"
)

func (s *Shell) Run() int {
	// for {
	fmt.Fprint(s.out, "$ ")

	line, err := s.reader.ReadString('\n')
	if err != nil {
		fmt.Fprint(s.out, err.Error())
		return 1
	}

	// removing the \n delimiter from the actual command so it doesn't mess up std output
	command := strings.Trim(line, "\n")

	_, ok := s.builtins[command]
	if !ok {
		fmt.Fprintf(s.out, "%s: command not found\n", command)
	}
	// }

	return 0
}
