package main

import (
	"bufio"
	"errors"
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
			// FIXME this sucks; figure out a better way to process aggregate errors
			fmt.Println(err)
		}
		return errors.New("multiple errors")

	}

	p := parser.New(tokens)
	expr, err := p.Parse()
	if err != nil {
		return err
	}

	i := interpreter.New(os.Stdout)
	if err := i.Interpret(expr); err != nil {
		return err
	}

	return nil
}

func fromFile(filename string) {
	input, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := execute(input); err != nil {
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
