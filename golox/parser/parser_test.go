package parser

import (
	"errors"
	"testing"

	"github.com/doeg/golox/golox/ast"
	"github.com/doeg/golox/golox/scanner"
	"github.com/doeg/golox/golox/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		testName      string
		input         string
		expected      ast.Expr
		expectedError error
	}{
		{
			input:    "123",
			expected: &ast.LiteralExpr{Value: float64(123)},
		},
		{
			input:    "\"hello, world\"",
			expected: &ast.LiteralExpr{Value: "hello, world"},
		},
		{
			input:    "true",
			expected: &ast.LiteralExpr{Value: true},
		},
		{
			input:    "false",
			expected: &ast.LiteralExpr{Value: false},
		},
		{
			input:    "nil",
			expected: &ast.LiteralExpr{Value: nil},
		},
		{
			input: "-1",
			expected: &ast.UnaryExpr{
				Operator: &token.Token{
					Lexeme: "-",
					Line:   0,
					Type:   token.MINUS,
				},
				Right: &ast.LiteralExpr{Value: float64(1)},
			},
		},
		{
			input: "5 - 3 - 1",
			expected: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: float64(5),
					},
					Operator: &token.Token{
						Lexeme: "-",
						Line:   0,
						Type:   token.MINUS,
					},
					Right: &ast.LiteralExpr{
						Value: float64(3),
					},
				},
				Operator: &token.Token{
					Lexeme: "-",
					Line:   0,
					Type:   token.MINUS,
				},
				Right: &ast.LiteralExpr{
					Value: float64(1),
				},
			},
		},
		{
			input: "5 - 3 * 1",
			expected: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: float64(5),
				},
				Operator: &token.Token{
					Lexeme: "-",
					Line:   0,
					Type:   token.MINUS,
				},
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: float64(3),
					},
					Operator: &token.Token{
						Lexeme: "*",
						Line:   0,
						Type:   token.STAR,
					},
					Right: &ast.LiteralExpr{
						Value: float64(1),
					},
				},
			},
		},
		{
			input: "(5 - 3) * 1",
			expected: &ast.BinaryExpr{
				Left: &ast.GroupingExpr{
					Expression: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: float64(5),
						},
						Operator: &token.Token{
							Lexeme: "-",
							Line:   0,
							Type:   token.MINUS,
						},
						Right: &ast.LiteralExpr{
							Value: float64(3),
						},
					},
				},
				Operator: &token.Token{
					Lexeme: "*",
					Line:   0,
					Type:   token.STAR,
				},
				Right: &ast.LiteralExpr{
					Value: float64(1),
				},
			},
		},
		{
			input: "-69 < 420 == true",
			expected: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.UnaryExpr{
						Operator: &token.Token{
							Lexeme: "-",
							Line:   0,
							Type:   token.MINUS,
						},
						Right: &ast.LiteralExpr{Value: float64(69)},
					},
					Operator: &token.Token{
						Lexeme: "<",
						Line:   0,
						Type:   token.LESS,
					},
					Right: &ast.LiteralExpr{
						Value: float64(420),
					},
				},
				Operator: &token.Token{
					Lexeme: "==",
					Line:   0,
					Type:   token.EQUAL_EQUAL,
				},
				Right: &ast.LiteralExpr{Value: true},
			},
		},
		{
			testName:      "error: expected expression",
			input:         "* 1",
			expectedError: errors.New(ErrExpectExpression),
		},
		{
			testName:      "error: missing closing paren",
			input:         "(1",
			expectedError: errors.New(ErrExpectClosingParen),
		},
	}

	for _, tt := range tests {
		testName := tt.testName
		if tt.testName == "" {
			testName = tt.input
		}

		tt := tt
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			s := scanner.New([]byte(tt.input))
			tokens, errors := s.ScanTokens()
			require.Empty(t, errors)

			p := New(tokens)
			expr, err := p.parseExpression()

			if tt.expectedError != nil {
				assert.Nil(t, expr)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.EqualValues(t, tt.expected, expr)
				assert.Nil(t, err)
			}
		})
	}
}

func TestSynchronize(t *testing.T) {
	tests := []struct {
		input         string
		expectedIndex int
	}{
		{
			input:         "(1 + 2; 3;",
			expectedIndex: 5,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			s := scanner.New([]byte(tt.input))
			tokens, errors := s.ScanTokens()
			require.Empty(t, errors)

			p := New(tokens)
			err := p.synchronize()
			require.Nil(t, err)

			assert.Equal(t, tt.expectedIndex, p.current)
		})
	}
}
