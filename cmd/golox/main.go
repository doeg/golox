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
	switch len(os.Args) {
	case 1:
		repl()
	case 2:
		fromFile(os.Args[1])
	default:
		fmt.Println("Usage: golox [filename]")
	}
}

func execute(input []byte) error {
	s := scanner.New([]byte(input))
	tokens, errs := s.ScanTokens()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
		return nil
	}

	p := parser.New(tokens)
	stmts, err := p.Parse()
	if err != nil {
		return err
	}

	i := interpreter.New()

	// TODO consider automatically printing the value of evaluated expressions...?
	_, err = i.Interpret(stmts)
	return err
}

func fromFile(filename string) {
	fmt.Println(filename)

	input, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = execute(input); err != nil {
		panic(err)
	}
}

func repl() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		if err := execute([]byte(text)); err != nil {
			fmt.Println("error: ", err)
			continue
		}
	}
}
