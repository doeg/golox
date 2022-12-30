package ast

import (
	"fmt"
	"strings"
)

type ASTPrinter struct{}

func (p *ASTPrinter) Print(expr Expr) string {
	result := expr.Accept(p).(string)
	return result
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	strs := make([]string, 0)
	for _, expr := range exprs {
		strs = append(strs, expr.Accept(p).(string))
	}

	return fmt.Sprintf("(%s %s)", name, strings.Join(strs, " "))
}

func (p *ASTPrinter) VisitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%+v", expr.Value)
}

func (p *ASTPrinter) VisitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}
