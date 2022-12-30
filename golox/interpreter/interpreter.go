package interpreter

import "github.com/doeg/golox/golox/ast"

type Interpreter struct{}

func New() *Interpreter {
	return &Interpreter{}
}

// evaluate sends the expression back into the interpreter's
// visitor implementation
func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.accept(i)
}

func (i *Interpreter) visitBinaryExpr(expr *ast.Binary) any {
	return nil
}

func (i *Interpreter) visitGroupingExpr(expr *ast.Grouping) any {
	return nil
}

func (i *Interpreter) visitUnaryExpr(expr *ast.Unary) any {
	return nil
}

func (i *Interpreter) visitLiteralExpr(expr *ast.Literal) any {
	return expr.Value
}
