package reader

import (
	"io"
	"strings"

	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),
)

// bellCompleter is a wrapper around the readline's AutoCompleter.
// all is does is copy the same structure as the underlying Do method and pass input to the Completer and see if there's any possible completion returned.
// if so, return them to the readline instance.
// if not, play the bell sound by writing the byte code for it
type bellCompleter struct {
	inner readline.AutoCompleter
	out   io.Writer
}

func (b *bellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// the current line only contains spaces with no characters typed so there is no possible completions + we should not play the bell sound.
	//
	// note: pos is where the cursor is currently at and we should only remove spaces before that as that is where the "line" is considered.
	// anything after the cursor is not part of the auto complete anymore.
	// ex. $     			█echo
	//
	// where █ is the current position of the cursor.
	// if you hit <Tab> from that position, there is no completions to be enabled because the line is starts for $ and ends at the cursor.
	//
	// ex. $ e█         echo
	// if you hit <Tab> here, then you would get the recommendations of echo.
	// to support this, we will only trim the space upto the cursor position and if there
	// are characters there then we will check for completions.
	if strings.TrimSpace(string(line[:pos])) == "" {
		return nil, 0
	}

	candidates, length := b.inner.Do(line, pos)
	if len(candidates) == 0 {
		// See: https://en.wikipedia.org/wiki/Bell_character
		// "\a" = bell sound
		io.WriteString(b.out, "\a")
	}
	return candidates, length
}
