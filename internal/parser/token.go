package parser

type TokenType string

const (
	Word     TokenType = "word"
	Operator TokenType = "operator"
)

type Token struct {
	Value string
	Type TokenType
}
