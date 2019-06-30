package main

import (
	"fmt"
)

func main() {
	code := "21 + 2 * (6 - 8 / 2)"
	lexer := NewLexer(code)
	parser := NewParser(lexer)
	ast, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	res, err := Eval(ast)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
