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
	doubleQuoteActive := false
	backslashActive := false

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
		// backslash check (\)
		// need to escape it first
		// handles the first backslash to turn the next character into a literal regardless of what it is
		if r == '\\' && backslashActive == false {
			backslashActive = true
			continue
		}

		// single quote check (')
		// if we encounter a singleQuote before we encounter a doubleQuote, then the singleQuote will be the wrapper that
		// has programmatic meaning while the doubleQuote is treated as a literal and gets added to the current token builder.
		// always setting hasToken to true so we don't flush and append an empty string as a token to final tokens output
		if r == '\'' && !doubleQuoteActive && !backslashActive {
			// toggle that sets to active on opening singleQuote and inactive on closing singleQuote so stringBuilder can build the token
			// and flush it to final output
			singleQuoteActive = !singleQuoteActive
			hasToken = true
			continue
		}
		// same as above but just for doubleQuote (")
		if r == '"' && !singleQuoteActive && !backslashActive {
			doubleQuoteActive = !doubleQuoteActive
			hasToken = true
			continue
		}

		// only flush if the singleQuote, doubleQuote, backSlash are not active, else the token is still being built
		if !singleQuoteActive && !doubleQuoteActive && !backslashActive && unicode.IsSpace(r) {
			flush()
			continue
		}

		current.WriteRune(r)
		hasToken = true
		backslashActive = false // after treating the character after the \ as a literal, reset backslash to false
	}

	// flush one final time to capture the last token if there is one and not <space><space><space>
	// at the end of the line
	flush()
	return tokens
}
