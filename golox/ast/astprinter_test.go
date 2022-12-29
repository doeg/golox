package ast

import (
	"testing"

	"github.com/doeg/golox/golox/token"
	"github.com/stretchr/testify/assert"
)

func TestASTPrinter(t *testing.T) {
	tests := []struct {
		input    Expr
		expected string
	}{
		{
			input: &Binary{
				Left: &Unary{
					Operator: &token.Token{
						Lexeme: "-",
						Line:   1,
						Type:   token.MINUS,
					},
					Right: &Literal{
						Value: 123,
					},
				},
				Operator: &token.Token{
					Lexeme: "*",
					Line:   1,
					Type:   token.STAR,
				},
				Right: &Grouping{
					Expression: &Literal{
						Value: 45.67,
					},
				},
			},
			expected: "(* (- 123) (group 45.67))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			p := &ASTPrinter{}
			result := p.Print(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
