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
	singleQuoteActive := false

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
		// single quote ' check
		// ' is not a valid argument so it is not appended to the string builder as part of a token.
		// it only exists to extend tokens by allowing spaces to be considered a token to be passed into commands.
		// adding continue to both opening & closing singleQuote check to ensure they are not written
		// to token builder
		if r == '\'' && !singleQuoteActive { // opening single quote check
			hasToken = true
			singleQuoteActive = true
			continue
		} else if r == '\'' && singleQuoteActive { // closing single quote check
			singleQuoteActive = false
			continue
		}
		// only flush if the singleQuote is not active, else the token is still being built
		if !singleQuoteActive && unicode.IsSpace(r) {
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
