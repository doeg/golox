package scanner

import (
	"fmt"
	"testing"

	"github.com/doeg/golox/golox/loxerror"
	"github.com/doeg/golox/golox/token"
	"github.com/stretchr/testify/assert"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		testName       string
		input          string
		expected       []token.Token
		expectedErrors []loxerror.LoxError
	}{
		{
			input: "( )",
			expected: []token.Token{
				{Line: 0, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 0, Lexeme: ")", Type: token.RIGHT_PAREN},
			},
		},
		{
			input: "\"string literal\"",
			expected: []token.Token{
				{Line: 0, Lexeme: "\"string literal\"", Literal: "string literal", Type: token.STRING},
			},
		},
		{
			input: "\"multiline\nstring\"",
			expected: []token.Token{
				{Line: 1, Lexeme: "\"multiline\nstring\"", Literal: "multiline\nstring", Type: token.STRING},
			},
		},
		{
			testName: "unterminated string",
			input:    "\"",
			expectedErrors: []loxerror.LoxError{
				{Line: 0, Message: loxerror.ErrUnterminatedString},
			},
		},
		{
			input: "1234",
			expected: []token.Token{
				{Line: 0, Lexeme: "1234", Literal: float64(1234), Type: token.NUMBER},
			},
		},
		{
			input: "12.34",
			expected: []token.Token{
				{Line: 0, Lexeme: "12.34", Literal: 12.34, Type: token.NUMBER},
			},
		},
		{
			input: "1 != 2",
			expected: []token.Token{
				{Line: 0, Lexeme: "1", Literal: float64(1), Type: token.NUMBER},
				{Line: 0, Lexeme: "!=", Type: token.BANG_EQUAL},
				{Line: 0, Lexeme: "2", Literal: float64(2), Type: token.NUMBER},
			},
		},
		{
			input: "// this is a comment\n!=",
			expected: []token.Token{
				{Line: 1, Lexeme: "!=", Type: token.BANG_EQUAL},
			},
		},
		{
			testName: "invalid character",
			input:    "@",
			expected: nil,
			expectedErrors: []loxerror.LoxError{
				{Line: 0, Message: fmt.Sprintf(loxerror.ErrUnexpectedCharacter, '@')},
			},
		},
		{
			testName: "invalid characters (multiline)",
			input:    "@\n$",
			expected: nil,
			expectedErrors: []loxerror.LoxError{
				{Line: 0, Message: fmt.Sprintf(loxerror.ErrUnexpectedCharacter, '@')},
				{Line: 1, Message: fmt.Sprintf(loxerror.ErrUnexpectedCharacter, '$')},
			},
		},
		{
			input: "var language = \"lox\";",
			expected: []token.Token{
				{Line: 0, Lexeme: "var", Type: token.VAR},
				{Line: 0, Lexeme: "language", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "=", Type: token.EQUAL},
				{Line: 0, Lexeme: "\"lox\"", Literal: "lox", Type: token.STRING},
				{Line: 0, Lexeme: ";", Type: token.SEMICOLON},
			},
		},
	}

	for _, tt := range tests {
		testName := tt.testName
		if tt.testName == "" {
			testName = tt.input
		}

		t.Run(testName, func(t *testing.T) {
			scanner := New([]byte(tt.input))
			tokens, errors := scanner.ScanTokens()

			switch {
			case len(tt.expectedErrors) > 0:
				assert.EqualValues(t, tt.expectedErrors, errors)
			default:
				assert.EqualValues(t, tt.expected, tokens)
				assert.Empty(t, errors)
			}
		})
	}
}
