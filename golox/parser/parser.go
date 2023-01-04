package parser

import (
	"errors"
	"fmt"

	"github.com/doeg/golox/golox/ast"
	"github.com/doeg/golox/golox/token"
)

var (
	ErrExpectClosingParen = "expect ')' after expression"
	ErrExpectExpression   = "expect expression"
)

// Parser implements Lox's grammar rules as a collection of methods.
// Each method for parsing a grammar rule produces a syntax tree for that
// rule and returns it to the caller.
type Parser struct {
	tokens []*token.Token

	// current points to the next token to be parsed
	current int
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		current: 0,
		tokens:  tokens,
	}
}

// Parse parses as many statements as we find until we reach EOF.
// This is equivalent to the grammar rule:
//
//	program -> declaration* EOF ;
func (p *Parser) Parse() ([]ast.Stmt, error) {
	statements := make([]ast.Stmt, 0)

	// TODO the book returns nil here :thinking:
	for {
		atEnd, err := p.isAtEnd()
		if err != nil {
			// TODO return Lox parse error
			return nil, err
		} else if atEnd {
			break
		}

		// Again, this is equivalent to the grammar rule:
		// 	program -> declaration* EOF ;
		//
		// In other words, a Lox program is a series of declarations.
		stmt, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

// advance consumes the current token and returns it
func (p *Parser) advance() (*token.Token, error) {
	done, err := p.isAtEnd()
	if err != nil {
		return nil, err
	}

	if !done {
		p.current++
	}

	return p.previous()
}

// check returns true if the current token is of the given type.
// Unlike match, it never consumes the token, it only looks at it.
func (p *Parser) check(tokenType token.TokenType) (bool, error) {
	done, err := p.isAtEnd()
	if err != nil {
		return false, err
	}

	if done {
		return false, nil
	}

	nextToken, err := p.peek()
	if err != nil {
		return false, err
	}

	return nextToken.Type == tokenType, nil
}

// consume checks to see if the next token is of the expected type.
// If so, it consumes the token. If some other token is there, then we've
// hit an error.
func (p *Parser) consume(tokenType token.TokenType, message string) (*token.Token, error) {
	isMatch, err := p.check(tokenType)
	if err != nil {
		return nil, err
	}

	if !isMatch {
		return nil, errors.New(message)
	}

	tok, err := p.advance()
	if err != nil {
		return nil, err
	}

	return tok, nil
}

// get returns a pointer to the Token at the given index.
func (p *Parser) get(index int) (*token.Token, error) {
	if index < 0 || index >= len(p.tokens) {
		return nil, fmt.Errorf("index %d out of bounds", index)
	}

	return p.tokens[index], nil
}

func (p *Parser) isAtEnd() (bool, error) {
	nextToken, err := p.peek()
	if err != nil {
		return false, err
	}

	return nextToken.Type == token.EOF, nil
}

// match checks to see if the current token has any of the given types.
// If so, it consumes the token and returns `true`. Otherwise, it returns `false`
// and leaves the current token alone.
func (p *Parser) match(tokenTypes ...token.TokenType) (bool, error) {
	for _, tokenType := range tokenTypes {
		isMatch, err := p.check(tokenType)
		if err != nil {
			return false, err
		}

		if isMatch {
			p.advance()
			return true, nil
		}
	}
	return false, nil
}

// parseComparison implements the following grammar rule:
//
//	comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term)* ;
func (p *Parser) parseComparison() (ast.Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		isMatch, err := p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
		if err != nil {
			return nil, err
		}

		if !isMatch {
			break
		}

		operator, err := p.previous()
		if err != nil {
			return nil, err
		}

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// parseDeclaration implements the following grammar rule:
//
//	declaration -> varDecl | statement ;
func (p *Parser) parseDeclaration() (ast.Stmt, error) {
	// TODO
	//
	// if the next token is a VAR, then parse and return a VarStmt by way of calling
	// p.parseVarDeclaration
	//
	// otherwise, return p.parseStatement
	return nil, nil
}

// parseEquality implements the following grammar rule:
//
//	equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) parseEquality() (ast.Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for {
		isMatch, err := p.match(token.BANG_EQUAL, token.EQUAL_EQUAL)
		if err != nil {
			return nil, err
		}

		if !isMatch {
			break
		}

		operator, err := p.previous()
		if err != nil {
			return nil, err
		}

		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// parseExpression implements the following grammar rule:
//
//	expression -> equality ;
func (p *Parser) parseExpression() (ast.Expr, error) {
	return p.parseEquality()
}

// parseExpressionStmt implements the following grammar rule:
//
//	exprStmt -> expression ";" ;
func (p *Parser) parseExpressionStmt() (ast.Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(token.SEMICOLON, "expected ';' after value"); err != nil {
		return nil, err
	}

	return &ast.ExpressionStmt{
		Expression: expr,
	}, nil
}

// parseFactor implements the following grammar rule:
//
//	factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) parseFactor() (ast.Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		isMatch, err := p.match(token.SLASH, token.STAR)
		if err != nil {
			return nil, err
		}

		if !isMatch {
			break
		}

		operator, err := p.previous()
		if err != nil {
			return nil, err
		}

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// parsePrimary implements the following grammar rule:
//
//	primary -> 	NUMBER | STRING | "true" | "false" | "nil"
//				| "(" expression ")"
//				| IDENTIFIER ;
func (p *Parser) parsePrimary() (ast.Expr, error) {
	isMatch, err := p.match(token.FALSE)
	if err != nil {
		return nil, err
	} else if isMatch {
		return &ast.LiteralExpr{Value: false}, nil
	}

	isMatch, err = p.match(token.TRUE)
	if err != nil {
		return nil, err
	} else if isMatch {
		return &ast.LiteralExpr{Value: true}, nil
	}

	isMatch, err = p.match(token.NIL)
	if err != nil {
		return nil, err
	} else if isMatch {
		return &ast.LiteralExpr{Value: nil}, nil
	}

	isMatch, err = p.match(token.NUMBER, token.STRING)
	if err != nil {
		return nil, err
	} else if isMatch {
		prev, err := p.previous()
		if err != nil {
			return nil, err
		}

		return &ast.LiteralExpr{Value: prev.Literal}, err
	}

	isMatch, err = p.match(token.LEFT_PAREN)
	if err != nil {
		return nil, err
	} else if isMatch {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(token.RIGHT_PAREN, ErrExpectClosingParen); err != nil {
			return nil, err
		}

		return &ast.GroupingExpr{
			Expression: expr,
		}, nil
	}

	// TODO return a LoxError instead of a regular error for unrecognized type
	return nil, errors.New(ErrExpectExpression)
}

// parsePrint implements the following grammar rule:
//
//	printStmt -> "print" expression ";" ;
func (p *Parser) parsePrint() (ast.Stmt, error) {
	// TODO
	return nil, nil
}

// parseStatement implements the following grammar rule:
//
//	statement -> exprStmt | printStmt ;
func (p *Parser) parseStatement() (ast.Stmt, error) {
	// TODO
	return nil, nil
}

// parseTerm implements the following grammar rule:
//
//	term -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) parseTerm() (ast.Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for {
		isMatch, err := p.match(token.MINUS, token.PLUS)
		if err != nil {
			return nil, err
		}

		if !isMatch {
			break
		}

		operator, err := p.previous()
		if err != nil {
			return nil, err
		}

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// parseUnary implements the following grammar rule:
//
//	unary -> ( "!" | "-" ) unary
//		 	 | primary
func (p *Parser) parseUnary() (ast.Expr, error) {
	isMatch, err := p.match(token.BANG, token.MINUS)
	if err != nil {
		return nil, err
	}

	if isMatch {
		operator, err := p.previous()
		if err != nil {
			return nil, err
		}

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.parsePrimary()
}

// This implements the following grammar rule:
//
//	varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) parseVarDeclaration() (ast.Stmt, error) {
	// TODO
	return nil, nil
}

// peek is a one-token lookahead, returning the current token without consuming it.
func (p *Parser) peek() (*token.Token, error) {
	return p.get(p.current)
}

func (p *Parser) previous() (*token.Token, error) {
	return p.get(p.current - 1)
}

// synchronize discards tokens until we're at the beginning of a new statement.
func (p *Parser) synchronize() error {
	_, err := p.advance()
	if err != nil {
		return err
	}

	for {
		atEnd, err := p.isAtEnd()
		if err != nil {
			return err
		}

		if atEnd {
			return nil
		}

		prev, err := p.previous()
		if err != nil {
			return err
		}

		if prev.Type == token.SEMICOLON {
			return nil
		}

		nextToken, err := p.peek()
		if err != nil {
			return err
		}

		switch nextToken.Type {
		case token.CLASS, token.FOR, token.FUN, token.IF, token.PRINT, token.RETURN, token.VAR, token.WHILE:
			return nil
		}

		_, err = p.advance()
		if err != nil {
			return err
		}
	}
}
