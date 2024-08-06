package main

import (
	"flag"
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
	flags := struct {
		graph *string
		out   *string

		run *string
	}{
		graph: flag.String("graph", "", "produce graph"),
		out:   flag.String("out", "./ast.gv", "output file for graph"),

		run: flag.String("run", "", "execute file"),
	}

	flag.Parse()

	if flags.graph != nil && *flags.graph != "" {
		doGraph(*flags.graph, *flags.out)
	} else if flags.run != nil && *flags.run != "" {
		doRunFile(*flags.run)
	} else {
		doREPL()
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

func doGraph(src string, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error opening file: %v\n", err)
		return
	}

	graph.Graph(
		string(data),
		dst,
	)
}
