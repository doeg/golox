package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/doeg/golox/golox/interpreter"
	"github.com/doeg/golox/golox/parser"
	"github.com/doeg/golox/golox/scanner"
)

func main() {
	repl()
}

func repl() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		s := scanner.New([]byte(text))
		tokens, errs := s.ScanTokens()
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Println(err.Error())
			}
			continue
		}

		p := parser.New(tokens)
		expr, err := p.Parse()
		if err != nil {
			fmt.Println(err)
			continue
		}

		i := interpreter.New()
		result, err := i.Interpret(expr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%+v\n", result)
	}
}
