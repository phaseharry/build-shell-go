package reader

import (
	"io"
	"strings"

	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),
	readline.PcItem("type"),
	readline.PcItem("pwd"),
	readline.PcItem("cd"),
)

// bellCompleter is a wrapper around the readline's AutoCompleter.
// all it does is copy the same structure as the underlying Do method and pass input to the Completer and see if there's any possible completion returned.
// if so, return them to the readline instance.
// if not, play the bell sound by writing the byte code for it - except when there was no
// prefix to complete in the first place, which returns early and stays silent.
type bellCompleter struct {
	inner readline.AutoCompleter
	out   io.Writer
}

func (b *bellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// only the text before the cursor counts as a completion prefix - the
	// underlying completer does the same slice (complete_helper.go:113 takes
	// line[:pos] then TrimSpaceLeft), so anything after the cursor is ignored.
	//
	// when that prefix is blank we deliberately suppress completion. note this
	// is NOT "no matches": a blank prefix matches every builtin, so readline
	// would drop into complete mode and print the whole menu. we return nothing
	// instead, and stay silent rather than bell - a bell would be wrong for a
	// case that has too many matches rather than none.
	//
	// ex. "   █echo"   -> prefix "   " -> blank -> suppressed
	//     "e█    echo" -> prefix "e"   -> matches echo AND exit -> readline
	//                     displays both
	//     "c█"         -> prefix "c"   -> only cd matches -> inserts "d "
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
