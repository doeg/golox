package loxerror

import "fmt"

type LoxError struct {
	Line    int
	Message string
}

var (
	ErrUnexpectedCharacter = "unexpected character %x"
	ErrUnterminatedString  = "unterminated string"
)

func (e *LoxError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s\n", e.Line, e.Message)
}
