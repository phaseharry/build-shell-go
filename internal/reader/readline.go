package reader

import (
	"github.com/chzyer/readline"
)

func New(cfg Config) Reader {
	readlineInstance, err := readline.NewEx(
		&readline.Config{
			Prompt: cfg.Prompt,
			AutoComplete: &bellCompleter{
				inner: completer,
				out:   cfg.Stdout,
			},
			InterruptPrompt: "^C",
			Stdin:           cfg.Stdin,
			Stdout:          cfg.Stdout,
			Stderr:          cfg.Stderr,
		})
	if err != nil {
		panic(err)
	}

	return readlineInstance
}
