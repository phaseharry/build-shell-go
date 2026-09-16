package reader

import (
	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),
)

func New(cfg Config) Reader {
	readlineInstance, err := readline.NewEx(&readline.Config{
		Prompt:          cfg.Prompt,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}

	return readlineInstance
}
