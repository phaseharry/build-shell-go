package parser

import (
	"strings"
	"unicode"
)

func ParseLine(line string) Command {
	command := Command{}
	tokens := tokenize(line)

	if len(tokens) == 0 {
		return command
	}

	commandToken := tokens[0]
	command.Name = commandToken.Value

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if token.Type == Word {
			command.Args = append(command.Args, token.Value)
			continue
		}

		// an operator always take the token directly after it as an argument.
		// ex. >> output.txt
		// if there is no token after it, then it is invalid and there is no destination to redirect an output to and should error out.
		// just ignoring it for now
		if i+1 > len(tokens) {
			break
		}

		operator := token.Value
		targetOperand := tokens[i+1].Value
		i++

		{
			command.Redirects = append(
				command.Redirects,
				Redirect{
					FileDescriptor: getFileDescriptor(operator),
					Target:         targetOperand,
					Append:         strings.HasSuffix(operator, ">>"),
				},
			)
		}
	}

	return command
}

func tokenize(line string) []Token {
	var tokens []Token
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
			currentString := current.String()
			token := Token{
				Value: currentString,
				Type:  Word,
			}
			tokens = append(tokens, token)
			current.Reset()
			hasToken = false
		}
	}

	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

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

		// redirection operator check (>)
		// only an unquoted, unescaped > is an operator.
		// the backslash block and both quote branches have already continued by this point, so whatever arrives here is unescaped.
		// we just need to make sure that neither quotes are active so the > is not being treated as a literal.
		if r == '>' && !singleQuoteActive && !doubleQuoteActive {
			operator := ">"

			// checking if the next character is > so it's the append operator >> instead
			// of the overwrite operator >
			if i+1 < len(runes) && runes[i+1] == '>' {
				operator = ">>"
				i += 1 // manually increment the index as we consume it here
			}
			// an unquoted run of digits touching the > is a file descriptor prefix (1>)
			// rather than a word of its own
			// checking if the current token includes any digits as those are valid with the > operator telling what file descriptors to redirect the output.
			// ie. 1>, 2>, 3>, etc.
			if hasToken && isDigits(current.String()) {
				operator = current.String() + operator
				current.Reset()
				hasToken = false
			} else {
				// if there isn't any digits as part of current actively built token, then we need to flush it
				// so the existing token that happens before the redirect operator gets used by the command before it.
				// ex. echo hi>output.txt
				// token{echo} token{hi} token{>} token{output.txt}
				flush()
			}
			tokens = append(tokens, Token{Value: operator, Type: Operator})
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

func isDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
