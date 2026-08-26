package shell

import "fmt"

func (s *Shell) Run() int {
	// for {
	fmt.Fprint(s.out, "$ ")

	line, err := s.reader.ReadString('\n')
	if err != nil {
		fmt.Fprint(s.out, err.Error())
		return 1
	}

	_, ok := s.builtins[line]
	if !ok {
		fmt.Fprintf(s.out, "%s: command not found\n", line)
	}
	// }

	return 0
}
