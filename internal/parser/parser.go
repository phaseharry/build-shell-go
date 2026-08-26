package parser

import (
	"strings"
	"unicode"
)

func ParseLine(line string) Command {
	line = strings.TrimSpace(line)

	// get the index of the first space delimiter. all the characters leading up to that first space will be the command
	// the result will be the arguments for that command
	commandDelimiter := strings.IndexFunc(line, unicode.IsSpace)

	// the line was either empty and has no command or
	// it was a command only with no arguments
	if commandDelimiter < 0 {
		return Command{
			Name: line,
		}
	}

	return Command{
		Name: line[:commandDelimiter],
		Args: line[commandDelimiter+1:],
	}
}
