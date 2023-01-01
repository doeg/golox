package interpreter

import (
	"errors"

	"github.com/doeg/golox/golox/ast"
	"github.com/doeg/golox/golox/token"
)

type Interpreter struct{}

func New() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr ast.Expr) (any, error) {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.BANG_EQUAL:
		eq, err := i.isEqual(left, right)
		if err != nil {
			return eq, err
		}
		return !eq, err
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case token.GREATER:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li > ri, err
	case token.GREATER_EQUAL:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li >= ri, err
	case token.LESS:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li < ri, err
	case token.LESS_EQUAL:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li <= ri, err
	case token.MINUS:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li - ri, err
	case token.PLUS:
		// Note, here we check for err == nil, NOT != nil
		li, ri, err := i.checkNumberOperands(left, right)
		if err == nil {
			return li + ri, nil
		}

		ls, rs, err := i.checkStringOperands(left, right)
		if err == nil {
			return ls + rs, nil
		}

		return nil, errors.New("operators must be strings or numbers")
	case token.SLASH:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li / ri, err
	case token.STAR:
		li, ri, err := i.checkNumberOperands(left, right)
		if err != nil {
			return nil, err
		}
		return li * ri, err
	}

	return nil, errors.New("invalid binary operator")
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	// Unary expressions have a single sub-expression that we evaluate first.
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		f, ok := right.(float64)
		if !ok {
			// TODO handle this more tidily (as a runtime error)
			return nil, errors.New("invalid cast")
		}
		return -f, nil
	case token.BANG:
		return !i.isTruthy(right), nil
	}

	// Unreachable. TODO: return an error...?
	return nil, nil
}

func (i *Interpreter) checkNumberOperands(left, right any) (float64, float64, error) {
	li, lok := left.(float64)
	ri, rok := right.(float64)

	if !lok || !rok {
		// TODO a better error message
		return li, ri, errors.New("operands must be numbers")
	}

	return li, ri, nil
}

func (i *Interpreter) checkStringOperands(left, right any) (string, string, error) {
	li, lok := left.(string)
	ri, rok := right.(string)

	if !lok || !rok {
		// TODO a better error message
		return li, ri, errors.New("operands must be strings")
	}

	return li, ri, nil
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) isEqual(a, b any) (bool, error) {
	// The book uses Java's .equals() method here, which may or
	// may not have some subtle difference from go's == operator.
	// Unlike the comparison operators, the book also allows
	// inequality checks on literals other than numbers.
	//
	// I don't really want to dig through all the various implementations
	// of that method, so I think leaning on == is fine for now.
	return a == b, nil
}

func (i *Interpreter) isTruthy(val any) bool {
	switch v := val.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}