//go:build ignore

package main

import (
	"bytes"
	_ "embed"
	"flag"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var outputFlag = flag.String("output", "./golox/ast/ast.go", "a filepath for the generated output")

var ast = []string{
	"Binary		:	Expr left, Token operator, Expr right",
	"Grouping	:	Expr expression",
	"Literal	:	Object value",
	"Unary		:	Token operator, Expr right",
}

//go:embed ast.go.tmpl
var tmpl string

func main() {
	flag.Parse()

	abs, err := filepath.Abs(*outputFlag)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(filepath.Dir(abs), 0700)

	f, err := os.Create(abs)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err = defineAST(f, "Expr", ast); err != nil {
		panic(err)
	}
}

type TemplateParams struct {
	BaseInterface    string
	Types            []ExpressionDef
	VisitorFunctions []string
}

type ExpressionDef struct {
	StructName string
	Fields     []ExpressionField
}

type ExpressionField struct {
	Name string
	Type string
}

// defineAST prints AST type definitions given a list of Lox grammar rules.
func defineAST(f *os.File, baseName string, grammar []string) error {
	defs := make([]ExpressionDef, 0)
	visitorFunctions := make([]string, 0)

	// This is pretty gross but, as the book mentions, robustness
	// ain't a priority :cowboy:
	for _, e := range grammar {
		parts := strings.Split(e, ":")

		structName := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])

		defs = append(defs, defineStruct(structName, fields))
		visitorFunctions = append(visitorFunctions, structName)
	}

	t, err := template.New("golox-ast").Funcs(template.FuncMap{
		// strings.Title is deprecated but it only really matters for Unicode text,
		// which Lox is not. :)
		"ToTitle": strings.Title,
	}).Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, TemplateParams{
		BaseInterface:    baseName,
		Types:            defs,
		VisitorFunctions: visitorFunctions,
	}); err != nil {
		return err
	}

	// Nicely format the go source code
	p, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if _, err := f.Write(p); err != nil {
		return err
	}

	return nil
}

// defineStruct generates the struct name and fields for a node in the AST; for example:
//
//	defineStruct("Unary", "Token operator, Expr right")
//
// ...will return:
//
//	{
//		StructName: "Unary",
//		Fields: [
//			{ Name: "operator", Type: "*token.Token" },
//			{ Name: "right", Type: "Expr" },
//		],
//	}
func defineStruct(structName, fieldList string) ExpressionDef {
	fields := make([]ExpressionField, 0)

	for _, fp := range strings.Split(fieldList, ",") {
		p := strings.Split(strings.TrimSpace(fp), " ")

		// Crafting Interpreters defines field types in a Java-ish way.
		// For the sake of... consistency? I've decided to keep the AST definition
		// (in its text form) exactly as the book has it. It does, however, necessitate
		// this ugly little switch statement to exorcise the Java-isms.
		fieldType := p[0]
		switch fieldType {
		case "Object":
			fieldType = "interface{}"
		case "Token":
			fieldType = "*token.Token"
		}

		fields = append(fields, ExpressionField{
			Name: p[1],
			Type: fieldType,
		})
	}

	return ExpressionDef{
		StructName: structName,
		Fields:     fields,
	}
}
