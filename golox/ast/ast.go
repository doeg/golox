// Code generated by golox-ast. DO NOT EDIT.
package ast

import (
	"github.com/doeg/golox/golox/token"
)

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
}

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type BinaryExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator *token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnaryExpr(e)
}

type StmtVisitor interface {
	VisitExpressionStmt(expr *ExpressionStmt) (any, error)
	VisitPrintStmt(expr *PrintStmt) (any, error)
}

type Stmt interface {
	Accept(StmtVisitor) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (e *PrintStmt) Accept(v StmtVisitor) (any, error) {
	return v.VisitPrintStmt(e)
}
