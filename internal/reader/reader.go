package reader

type Config struct {
	Prompt string
}

type Reader interface {
	Readline() (string, error)
}
