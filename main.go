package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/manuelpepe/interpreter/graph"
	"github.com/manuelpepe/interpreter/repl"
)

const option = 2

func main() {
	switch option {
	case 1:
		doREPL()
	case 2:
		doGraph()
	default:
		panic("unknown option")
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
