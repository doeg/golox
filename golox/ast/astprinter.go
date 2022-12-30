package ast

import (
	"fmt"
	"strings"
)

type ASTPrinter struct{}

func (p *ASTPrinter) Print(expr Expr) string {
	result, _ := expr.Accept(p)
	return result.(string)
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	strs := make([]string, 0)
	for _, expr := range exprs {
		s, err := expr.Accept(p)
		if err != nil {
			return "", err
		}

		strs = append(strs, s.(string))
	}

	return fmt.Sprintf("(%s %s)", name, strings.Join(strs, " ")), nil
}

func (p *ASTPrinter) VisitBinaryExpr(expr *Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return nil, nil
	}

	return fmt.Sprintf("%+v", expr.Value), nil
}

func (p *ASTPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}
