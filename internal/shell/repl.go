package shell

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

// Run starts the read-evaluate-print loop and returns the shell's exit status.
func (s *Shell) Run() int {
	for {
		if _, err := fmt.Fprint(s.out, "$ "); err != nil {
			return 1
		}

		line, readErr := s.reader.ReadString('\n')
		if len(line) > 0 {
			result := s.ExecuteLine(strings.TrimRight(line, "\r\n"))
			if result.Exit {
				return result.Status
			}
		}

		if readErr == nil {
			continue
		}
		if errors.Is(readErr, io.EOF) {
			return s.lastStatus
		}

		fmt.Fprintln(s.errOut, readErr)
		return 1
	}
}

// ExecuteLine parses and executes one line of shell input.
func (s *Shell) ExecuteLine(line string) Result {
	command, ok, err := parser.ParseLine(line)
	if err != nil {
		fmt.Fprintf(s.errOut, "syntax error: %v\n", err)
		s.lastStatus = 2
		return Result{Status: s.lastStatus}
	}
	if !ok {
		return Result{Status: s.lastStatus}
	}

	result := s.execute(command)
	s.lastStatus = result.Status

	return result
}
