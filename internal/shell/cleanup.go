package shell

import "fmt"

func (s *Shell) Cleanup() {
	if err := s.reader.Close(); err != nil {
		fmt.Fprintf(s.errOut, "issue with closing reader: %s", err.Error())
	}
}
