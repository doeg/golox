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
		expected       []*token.Token
		expectedErrors []loxerror.LoxError
	}{
		{
			input: "( )",
			expected: []*token.Token{
				{Line: 0, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 0, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 0, Type: token.EOF},
			},
		},

		{
			input: "\"string literal\"",
			expected: []*token.Token{
				{Line: 0, Lexeme: "\"string literal\"", Literal: "string literal", Type: token.STRING},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "\"multiline\nstring\"",
			expected: []*token.Token{
				{Line: 1, Lexeme: "\"multiline\nstring\"", Literal: "multiline\nstring", Type: token.STRING},
				{Line: 1, Type: token.EOF},
			},
		},
		{
			input: "1234",
			expected: []*token.Token{
				{Line: 0, Lexeme: "1234", Literal: float64(1234), Type: token.NUMBER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "12.34",
			expected: []*token.Token{
				{Line: 0, Lexeme: "12.34", Literal: 12.34, Type: token.NUMBER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "1*2 >= 3*4",
			expected: []*token.Token{
				{Line: 0, Lexeme: "1", Literal: float64(1), Type: token.NUMBER},
				{Line: 0, Lexeme: "*", Type: token.STAR},
				{Line: 0, Lexeme: "2", Literal: float64(2), Type: token.NUMBER},
				{Line: 0, Lexeme: ">=", Type: token.GREATER_EQUAL},
				{Line: 0, Lexeme: "3", Literal: float64(3), Type: token.NUMBER},
				{Line: 0, Lexeme: "*", Type: token.STAR},
				{Line: 0, Lexeme: "4", Literal: float64(4), Type: token.NUMBER},

				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "-12/0.34",
			expected: []*token.Token{
				{Line: 0, Lexeme: "-", Type: token.MINUS},
				{Line: 0, Lexeme: "12", Literal: float64(12), Type: token.NUMBER},
				{Line: 0, Lexeme: "/", Type: token.SLASH},
				{Line: 0, Lexeme: "0.34", Literal: 0.34, Type: token.NUMBER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "1 != 2",
			expected: []*token.Token{
				{Line: 0, Lexeme: "1", Literal: float64(1), Type: token.NUMBER},
				{Line: 0, Lexeme: "!=", Type: token.BANG_EQUAL},
				{Line: 0, Lexeme: "2", Literal: float64(2), Type: token.NUMBER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "a == b",
			expected: []*token.Token{
				{Line: 0, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "==", Type: token.EQUAL_EQUAL},
				{Line: 0, Lexeme: "b", Type: token.IDENTIFIER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "a != b",
			expected: []*token.Token{
				{Line: 0, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "!=", Type: token.BANG_EQUAL},
				{Line: 0, Lexeme: "b", Type: token.IDENTIFIER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "a = !b",
			expected: []*token.Token{
				{Line: 0, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "=", Type: token.EQUAL},
				{Line: 0, Lexeme: "!", Type: token.BANG},
				{Line: 0, Lexeme: "b", Type: token.IDENTIFIER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "a > b",
			expected: []*token.Token{
				{Line: 0, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: ">", Type: token.GREATER},
				{Line: 0, Lexeme: "b", Type: token.IDENTIFIER},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "// this is a comment\n!=",
			expected: []*token.Token{
				{Line: 1, Lexeme: "!=", Type: token.BANG_EQUAL},
				{Line: 1, Type: token.EOF},
			},
		},
		{
			input: "var language = \"lox\";",
			expected: []*token.Token{
				{Line: 0, Lexeme: "var", Type: token.VAR},
				{Line: 0, Lexeme: "language", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "=", Type: token.EQUAL},
				{Line: 0, Lexeme: "\"lox\"", Literal: "lox", Type: token.STRING},
				{Line: 0, Lexeme: ";", Type: token.SEMICOLON},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: "var average = (min + max)/2;",
			expected: []*token.Token{
				{Line: 0, Lexeme: "var", Type: token.VAR},
				{Line: 0, Lexeme: "average", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "=", Type: token.EQUAL},
				{Line: 0, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 0, Lexeme: "min", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "+", Type: token.PLUS},
				{Line: 0, Lexeme: "max", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 0, Lexeme: "/", Type: token.SLASH},
				{Line: 0, Lexeme: "2", Literal: float64(2), Type: token.NUMBER},
				{Line: 0, Lexeme: ";", Type: token.SEMICOLON},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: `
				var a = 1;
				while (a <= 10) {
					print a;
				}
			`,
			expected: []*token.Token{
				{Line: 1, Lexeme: "var", Type: token.VAR},
				{Line: 1, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "=", Type: token.EQUAL},
				{Line: 1, Lexeme: "1", Literal: float64(1), Type: token.NUMBER},
				{Line: 1, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 2, Lexeme: "while", Type: token.WHILE},
				{Line: 2, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 2, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: "<=", Type: token.LESS_EQUAL},
				{Line: 2, Lexeme: "10", Literal: float64(10), Type: token.NUMBER},
				{Line: 2, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 2, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 3, Lexeme: "print", Type: token.PRINT},
				{Line: 3, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 3, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 4, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 5, Type: token.EOF},
			},
		},
		{
			input: `
				if (condition) {
					print "yes";
				} else {
					print "no";
				}
			`,
			expected: []*token.Token{
				{Line: 1, Lexeme: "if", Type: token.IF},
				{Line: 1, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 1, Lexeme: "condition", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 1, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 2, Lexeme: "print", Type: token.PRINT},
				{Line: 2, Lexeme: "\"yes\"", Literal: "yes", Type: token.STRING},
				{Line: 2, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 3, Lexeme: "}", Type: token.RIGHT_BRACE},
				{Line: 3, Lexeme: "else", Type: token.ELSE},
				{Line: 3, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 4, Lexeme: "print", Type: token.PRINT},
				{Line: 4, Lexeme: "\"no\"", Literal: "no", Type: token.STRING},
				{Line: 4, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 5, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 6, Type: token.EOF},
			},
		},
		{
			input: `
				for (var a = 1; a < 10; a = a + 1) {
					print a;
				}
			`,
			expected: []*token.Token{
				{Line: 1, Lexeme: "for", Type: token.FOR},
				{Line: 1, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 1, Lexeme: "var", Type: token.VAR},
				{Line: 1, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "=", Type: token.EQUAL},
				{Line: 1, Lexeme: "1", Literal: float64(1), Type: token.NUMBER},
				{Line: 1, Lexeme: ";", Type: token.SEMICOLON},
				{Line: 1, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "<", Type: token.LESS},
				{Line: 1, Lexeme: "10", Literal: float64(10), Type: token.NUMBER},
				{Line: 1, Lexeme: ";", Type: token.SEMICOLON},
				{Line: 1, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "=", Type: token.EQUAL},
				{Line: 1, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "+", Type: token.PLUS},
				{Line: 1, Lexeme: "1", Literal: float64(1), Type: token.NUMBER},
				{Line: 1, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 1, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 2, Lexeme: "print", Type: token.PRINT},
				{Line: 2, Lexeme: "a", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 3, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 4, Type: token.EOF},
			},
		},
		{
			input: "makeBreakfast(bagel, creamCheese, lox);",
			expected: []*token.Token{
				{Line: 0, Lexeme: "makeBreakfast", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 0, Lexeme: "bagel", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: ",", Type: token.COMMA},
				{Line: 0, Lexeme: "creamCheese", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: ",", Type: token.COMMA},
				{Line: 0, Lexeme: "lox", Type: token.IDENTIFIER},
				{Line: 0, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 0, Lexeme: ";", Type: token.SEMICOLON},
				{Line: 0, Type: token.EOF},
			},
		},
		{
			input: `
				fun outerFunction() {
					fun innerFunction(str) {
						print str;
					}

					innerFunction("hello");
				}
			`,
			expected: []*token.Token{
				{Line: 1, Lexeme: "fun", Type: token.FUN},
				{Line: 1, Lexeme: "outerFunction", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 1, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 1, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 2, Lexeme: "fun", Type: token.FUN},
				{Line: 2, Lexeme: "innerFunction", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 2, Lexeme: "str", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 2, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 3, Lexeme: "print", Type: token.PRINT},
				{Line: 3, Lexeme: "str", Type: token.IDENTIFIER},
				{Line: 3, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 4, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 6, Lexeme: "innerFunction", Type: token.IDENTIFIER},
				{Line: 6, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 6, Lexeme: "\"hello\"", Literal: "hello", Type: token.STRING},
				{Line: 6, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 6, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 7, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 8, Type: token.EOF},
			},
		},
		{
			input: `
				class Brunch < Breakfast {
					init(meat, bread, drink) {
						super.init(meat, bread);
						this.drink = drink;
					}
				}
			`,
			expected: []*token.Token{
				{Line: 1, Lexeme: "class", Type: token.CLASS},
				{Line: 1, Lexeme: "Brunch", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "<", Type: token.LESS},
				{Line: 1, Lexeme: "Breakfast", Type: token.IDENTIFIER},
				{Line: 1, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 2, Lexeme: "init", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 2, Lexeme: "meat", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: ",", Type: token.COMMA},
				{Line: 2, Lexeme: "bread", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: ",", Type: token.COMMA},
				{Line: 2, Lexeme: "drink", Type: token.IDENTIFIER},
				{Line: 2, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 2, Lexeme: "{", Type: token.LEFT_BRACE},

				{Line: 3, Lexeme: "super", Type: token.SUPER},
				{Line: 3, Lexeme: ".", Type: token.DOT},
				{Line: 3, Lexeme: "init", Type: token.IDENTIFIER},
				{Line: 3, Lexeme: "(", Type: token.LEFT_PAREN},
				{Line: 3, Lexeme: "meat", Type: token.IDENTIFIER},
				{Line: 3, Lexeme: ",", Type: token.COMMA},
				{Line: 3, Lexeme: "bread", Type: token.IDENTIFIER},
				{Line: 3, Lexeme: ")", Type: token.RIGHT_PAREN},
				{Line: 3, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 4, Lexeme: "this", Type: token.THIS},
				{Line: 4, Lexeme: ".", Type: token.DOT},
				{Line: 4, Lexeme: "drink", Type: token.IDENTIFIER},
				{Line: 4, Lexeme: "=", Type: token.EQUAL},
				{Line: 4, Lexeme: "drink", Type: token.IDENTIFIER},
				{Line: 4, Lexeme: ";", Type: token.SEMICOLON},

				{Line: 5, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 6, Lexeme: "}", Type: token.RIGHT_BRACE},

				{Line: 7, Type: token.EOF},
			},
		},
		{
			testName: "error: unterminated string",
			input:    "\"",
			expectedErrors: []loxerror.LoxError{
				{Line: 0, Message: loxerror.ErrUnterminatedString},
			},
		},
		{
			testName: "error: invalid character",
			input:    "@",
			expected: nil,
			expectedErrors: []loxerror.LoxError{
				{Line: 0, Message: fmt.Sprintf(loxerror.ErrUnexpectedCharacter, '@')},
			},
		},
		{
			testName: "error: invalid characters (multiline)",
			input:    "@\n$",
			expected: nil,
			expectedErrors: []loxerror.LoxError{
				{Line: 0, Message: fmt.Sprintf(loxerror.ErrUnexpectedCharacter, '@')},
				{Line: 1, Message: fmt.Sprintf(loxerror.ErrUnexpectedCharacter, '$')},
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
