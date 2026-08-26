package parser

import (
	"strings"
	"unicode"
)

func ParseLine(line string) Command {
	tokens := tokenize(line)

	if len(tokens) == 0 {
		return Command{}
	}

	commandName := tokens[0]
	var args []string

	if len(tokens) > 1 {
		args = tokens[1:]
	}

	return Command{
		Name: commandName,
		Args: args,
	}
}

func tokenize(line string) []string {
	var tokens []string
	var current strings.Builder
	hasToken := false

	flush := func() {
		// everytime flush is called, convert the current string builder to a string and reset it to build the next token.
		// flush is called everytime we encounter a <space> or if we are at the end of the line and we need to create the
		// very last string
		if hasToken {
			tokens = append(tokens, current.String())
			current.Reset()
			hasToken = false
		}
	}

	for _, r := range line {
		if unicode.IsSpace(r) {
			flush()
			continue
		}
		current.WriteRune(r)
		hasToken = true
	}

	// flush one final time to capture the last token if there is one and not <space><space><space>
	// at the end of the line
	flush()
	return tokens
}
