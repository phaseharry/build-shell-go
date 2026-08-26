package parser

import (
	"errors"
	"strings"
	"unicode"
)

// ErrUnterminatedSingleQuote indicates that a quoted word was not closed.
var ErrUnterminatedSingleQuote = errors.New("unterminated single quote")

// ParseLine splits a line into a command while preserving whitespace inside
// single-quoted words. The boolean is false when the line contains no command.
func ParseLine(line string) (Command, bool, error) {
	var (
		words         []string
		word          strings.Builder
		inSingleQuote bool
		wordStarted   bool
	)

	for _, character := range line {
		if inSingleQuote {
			if character == '\'' {
				inSingleQuote = false
			} else {
				word.WriteRune(character)
			}
			continue
		}

		switch {
		case character == '\'':
			inSingleQuote = true
			wordStarted = true
		case unicode.IsSpace(character):
			if wordStarted {
				words = append(words, word.String())
				word.Reset()
				wordStarted = false
			}
		default:
			word.WriteRune(character)
			wordStarted = true
		}
	}

	if inSingleQuote {
		return Command{}, false, ErrUnterminatedSingleQuote
	}

	if wordStarted {
		words = append(words, word.String())
	}

	if len(words) == 0 {
		return Command{}, false, nil
	}

	return Command{
		Name: words[0],
		Args: words[1:],
	}, true, nil
}
