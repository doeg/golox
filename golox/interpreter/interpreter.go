package interpreter

import (
	"github.com/doeg/golox/golox/ast"
	"github.com/doeg/golox/golox/token"
)

type Interpreter struct{}

func New() *Interpreter {
	return &Interpreter{}
}

// evaluate sends the expression back into the interpreter's
// visitor implementation
func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) any {
	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) any {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		ir, ok := right.(float64)
		if !ok {
			// TODO
			panic("oh no")
		}
		return -ir
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) any {
	return expr.Value
}
