package scanner

import (
	"fmt"
	"strconv"

	"github.com/doeg/golox/golox/loxerror"
	"github.com/doeg/golox/golox/token"
)

type Scanner struct {
	source []byte
	tokens []token.Token

	// start points to the first character in the lexeme being scanned
	start int

	// current points at the character currently being considered
	current int

	// line tracks which source line current is on so we can produce
	// tokens that know their location
	line int

	// errors accumulates syntax errors as the scanner progresses
	// so as many errors as possible can be collected in a single scan pass
	errors []loxerror.LoxError
}

func New(source []byte) *Scanner {
	return &Scanner{
		errors: make([]loxerror.LoxError, 0),
		source: source,
		tokens: make([]token.Token, 0),
	}
}

func (scanner *Scanner) ScanTokens() ([]token.Token, []loxerror.LoxError) {
	for !scanner.isAtEnd() {
		// We are at the beginning of the next lexeme
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, token.Token{Line: scanner.line, Type: token.EOF})

	return scanner.tokens, scanner.errors
}

// addOperatorToken adds an operator token (i.e., a token without a literal value)
// to the list of tokens.
func (scanner *Scanner) addOperatorToken(tokenType token.TokenType) {
	scanner.addToken(tokenType, nil)
}

func (scanner *Scanner) addToken(tokenType token.TokenType, literal interface{}) {
	text := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(scanner.tokens, token.Token{
		Lexeme:  string(text),
		Line:    scanner.line,
		Literal: literal,
		Type:    tokenType,
	})
}

// advance consumes the next ASCII character in the source file and returns it.
func (scanner *Scanner) advance() byte {
	b := scanner.source[scanner.current]
	scanner.current++
	return b
}

// advanceUntilNewline advances the scanner until a newline character is encountered,
// or the scanner reaches the end of 'source'. Note that the newline character
// is not consumed, nor is 'scanner.line' advanced.
func (scanner *Scanner) advanceUntilNewline() {
	for scanner.peek() != '\n' && !scanner.isAtEnd() {
		scanner.advance()
	}
}

// isAtEnd returns true when all of the characters in 'source' have been consumed.
func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

// match returns true and consumes the current character if it matches 'expected'.
// In other words, it is like a conditional 'advance()'.
func (scanner *Scanner) match(expected byte) bool {
	if !scanner.isAtEnd() && scanner.peek() == expected {
		scanner.current++
		return true
	}

	return false
}

// peek is a one-character lookahead, returning the current character without consuming it.
func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return '\x00'
	}
	return scanner.source[scanner.current]
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return '\x00'
	}

	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) recordError(message string) {
	scanner.errors = append(scanner.errors, loxerror.LoxError{
		Line:    scanner.line,
		Message: message,
	})
}

func (scanner *Scanner) scanIdentifier() {
	for isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.source[scanner.start:scanner.current]
	tokenType, ok := token.Keywords[string(text)]
	if ok {
		scanner.addOperatorToken(tokenType)
	} else {
		scanner.addOperatorToken(token.IDENTIFIER)
	}
}

func (scanner *Scanner) scanNumber() {
	for isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
		// Consume the decimal
		scanner.advance()

		for isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	str := scanner.source[scanner.start:scanner.current]
	dbl, _ := strconv.ParseFloat(string(str), 64)
	scanner.addToken(token.NUMBER, dbl)

}

func (scanner *Scanner) scanString() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		scanner.recordError(loxerror.ErrUnterminatedString)
		return
	}

	// Consume the closing '"' character
	scanner.advance()

	// Trim the surrounding quotes from the string literal
	val := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addToken(token.STRING, string(val))
}

func (scanner *Scanner) scanToken() {
	b := scanner.advance()
	switch b {
	case '(':
		scanner.addOperatorToken(token.LEFT_PAREN)
	case ')':
		scanner.addOperatorToken(token.RIGHT_PAREN)
	case '{':
		scanner.addOperatorToken(token.LEFT_BRACE)
	case '}':
		scanner.addOperatorToken(token.RIGHT_BRACE)
	case ',':
		scanner.addOperatorToken(token.COMMA)
	case '.':
		scanner.addOperatorToken(token.DOT)
	case '-':
		scanner.addOperatorToken(token.MINUS)
	case '+':
		scanner.addOperatorToken(token.PLUS)
	case ';':
		scanner.addOperatorToken(token.SEMICOLON)
	case '/':
		if scanner.match('/') {
			scanner.advanceUntilNewline()
		} else {
			scanner.addOperatorToken(token.SLASH)
		}
	case '*':
		scanner.addOperatorToken(token.STAR)
	case '!':
		if scanner.match('=') {
			scanner.addOperatorToken(token.BANG_EQUAL)
		} else {
			scanner.addOperatorToken(token.BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.addOperatorToken(token.EQUAL_EQUAL)
		} else {
			scanner.addOperatorToken(token.EQUAL)
		}
	case '>':
		if scanner.match('=') {
			scanner.addOperatorToken(token.GREATER_EQUAL)
		} else {
			scanner.addOperatorToken(token.GREATER)
		}
	case '<':
		if scanner.match('=') {
			scanner.addOperatorToken(token.LESS_EQUAL)
		} else {
			scanner.addOperatorToken(token.LESS)
		}
	case '"':
		scanner.scanString()
	case ' ', '\r', '\t':
		// Continue, ignoring whitespace
	case '\n':
		scanner.line++
	default:
		switch {
		case isDigit(b):
			scanner.scanNumber()
		case isAlpha(b):
			scanner.scanIdentifier()
		default:
			scanner.recordError(fmt.Sprintf(loxerror.ErrUnexpectedCharacter, b))
		}
	}
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isAlphaNumeric(b byte) bool {
	return isAlpha(b) || isDigit(b)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
