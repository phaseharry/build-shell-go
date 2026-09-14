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
		// - handles the first backslash so the character after it can be turned into a literal
		// - check if the backslash is enclosed by a singleQuote. if so, then it has no special meaning and treat it as a literal value
		if r == '\\' && !backslashActive && !singleQuoteActive {
			backslashActive = true
			continue
		}

		// consume the character that the backslash above was holding open for.
		// outside of quotes a backslash escapes anything, but inside a doubleQuote it only
		// escapes $ ` " and \. before any other character it has no special meaning, so the
		// backslash itself has to be written back out as a literal before the character.
		// we cannot know which case applies until we see the character, which is why the
		// decision is made here rather than back where backslashActive was set.
		if backslashActive {
			backslashActive = false
			// if the character directly after the backslash is not any of these: % ` \
			// then the backslash is just a literal written to the current token
			if doubleQuoteActive && !strings.ContainsRune("$`\"\\", r) {
				current.WriteRune('\\')
			}
			current.WriteRune(r)
			hasToken = true
			continue
		}

		// single quote check (')
		// if we encounter a singleQuote before we encounter a doubleQuote, then the singleQuote will be the wrapper that
		// has programmatic meaning while the doubleQuote is treated as a literal and gets added to the current token builder.
		// always setting hasToken to true so we don't flush and append an empty string as a token to final tokens output
		// an escaped quote never reaches here because the backslash block above consumes it,
		// so this can only ever see a quote that really is a delimiter
		if r == '\'' && !doubleQuoteActive {
			// toggle that sets to active on opening singleQuote and inactive on closing singleQuote so stringBuilder can build the token and flush it to final output
			singleQuoteActive = !singleQuoteActive
			hasToken = true
			continue
		}
		// same as above but just for doubleQuote (")
		if r == '"' && !singleQuoteActive {
			doubleQuoteActive = !doubleQuoteActive
			hasToken = true
			continue
		}

		// only flush if the singleQuote and doubleQuote are not active, else the token is still being built
		if !singleQuoteActive && !doubleQuoteActive && unicode.IsSpace(r) {
			flush()
			continue
		}

		current.WriteRune(r)
		hasToken = true
	}

	// flush one final time to capture the last token if there is one and not <space><space><space> at the end of the line
	flush()
	return tokens
}
