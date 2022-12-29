package ast

import (
	"fmt"
	"strings"
)

type ASTPrinter struct{}

func (p *ASTPrinter) Print(expr Expr) string {
	result := expr.accept(p).(string)
	return result
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	strs := make([]string, 0)
	for _, expr := range exprs {
		strs = append(strs, expr.accept(p).(string))
	}

	return fmt.Sprintf("(%s %s)", name, strings.Join(strs, " "))
}

func (p *ASTPrinter) visitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) visitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) visitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%+v", expr.Value)
}

func (p *ASTPrinter) visitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}
