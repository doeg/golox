package ast

import (
	"fmt"
	"strings"

	"github.com/doeg/golox/golox/ast"
)

type ASTPrinter struct{}

func (p *ASTPrinter) Print(expr ast.Expr) string {
	result, _ := expr.Accept(p)
	return result.(string)
}

func (p *ASTPrinter) parenthesize(name string, exprs ...ast.Expr) (string, error) {
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

func (p *ASTPrinter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	if expr.Value == nil {
		return nil, nil
	}

	return fmt.Sprintf("%+v", expr.Value), nil
}

func (p *ASTPrinter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}
