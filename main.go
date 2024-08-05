package main

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/manuelpepe/interpreter/eval"
	"github.com/manuelpepe/interpreter/graph"
	"github.com/manuelpepe/interpreter/lexer"
	"github.com/manuelpepe/interpreter/object"
	"github.com/manuelpepe/interpreter/parser"
	"github.com/manuelpepe/interpreter/repl"
)

const option = 3

func main() {
	switch option {
	case 1:
		doREPL()
	case 2:
		doGraph()
	case 3:
		doRunFile(os.Args[1])
	default:
		panic("unknown option")
	}
}

func doRunFile(srcFile string) {
	data, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error opening file: %v\n", err)
		return
	}

	l := lexer.NewLexer(string(data))
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(os.Stdout, p.Errors())
		return
	}

	env := object.NewEnvironment()
	res := eval.Eval(prog, env)
	if res != nil {
		io.WriteString(os.Stdout, res.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}

func doREPL() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}

func doGraph() {
	graph.Graph(
		"let a = fn(a,b,c) { return fn() { return 3 + 1 }() }; let b = a",
		"./ast.gv",
	)

}
