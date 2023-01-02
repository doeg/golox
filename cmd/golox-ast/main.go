//go:build ignore

package main

import (
	"bytes"
	_ "embed"
	"go/format"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Embed the ast.go template file so that we don't have to mess with
// file I/O, relative pathnames, etc. etc.
//
//go:embed ast.go.tmpl
var tmpl string

// The full path of the output file, relative to the project root
const outPath string = "./golox/ast/"
const outName string = "ast.go"

type TemplateParams struct {
	Interfaces []ASTInterface
}

type ASTInterface struct {
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

func main() {
	interfaces := []ASTInterface{
		defineAST("Expr", []string{
			"Binary		:	Expr left, Token operator, Expr right",
			"Grouping	:	Expr expression",
			"Literal	:	Object value",
			"Unary		:	Token operator, Expr right",
		}),
		defineAST("Stmt", []string{
			"Expression	: Expr expression",
			"Print		: Expr expression",
		}),
	}

	t, err := template.New("golox-ast").Funcs(template.FuncMap{
		"ToTitle": strings.Title,
	}).Parse(tmpl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, TemplateParams{
		Interfaces: interfaces,
	}); err != nil {
		panic(err)
	}

	// Nicely format the go source code :)
	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	abs, err := filepath.Abs(outPath)
	if err != nil {
		panic(err)
	}

	// The directory should pretty much _always_ exist,
	// but recursively create it just in case we're doing something dumb.
	os.MkdirAll(abs, 0700)

	// Create the ast.go file, overwriting anything that's currently there.
	filePath := path.Join(abs, outName)
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Finally, write the formatted ast.go source to disk.
	if _, err := f.Write(p); err != nil {
		panic(err)
	}
}

func defineAST(baseInterface string, rules []string) ASTInterface {
	defs := make([]ExpressionDef, 0)
	visitorFunctions := make([]string, 0)

	// This is pretty gross but, as the book mentions, robustness
	// ain't a priority :cowboy:
	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		structName := strings.TrimSpace(parts[0])
		fieldList := strings.Split(strings.TrimSpace(parts[1]), ",")

		fields := make([]ExpressionField, 0)
		for _, fp := range fieldList {
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

		defs = append(defs, ExpressionDef{
			StructName: structName,
			Fields:     fields,
		})
		visitorFunctions = append(visitorFunctions, structName)
	}

	return ASTInterface{
		BaseInterface:    baseInterface,
		Types:            defs,
		VisitorFunctions: visitorFunctions,
	}
}
