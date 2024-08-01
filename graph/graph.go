package graph

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"

	"github.com/manuelpepe/interpreter/ast"
	"github.com/manuelpepe/interpreter/lexer"
	"github.com/manuelpepe/interpreter/parser"
)

func Graph(src string, dst string) {
	l := lexer.NewLexer(src)
	p := parser.New(l)

	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("\t" + msg + "\n")
		}
	}

	g := doGraph(prog)

	file, _ := os.Create(dst)
	_ = draw.DOT(g, file)

	cmd := exec.Command("dot", "-Tsvg", "-O", dst)
	_ = cmd.Run()
}

func doGraph(prog *ast.Program) graph.Graph[int, int] {
	g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic(), graph.Rooted())
	offset := 0
	id := 0

	g.AddVertex(id, graph.VertexAttribute("label", reflect.TypeOf(prog).String()))
	id += 1
	offset += 1

	for _, s := range prog.Statements {
		id += offset
		offset = graphNode(g, s, 0, offset)
	}

	return g
}

func graphNode(g graph.Graph[int, int], node ast.Node, parent int, offset int) int {
	id := 1 + offset
	offset += 1

	nodestr := fmt.Sprintf("%s\n%s", reflect.TypeOf(node).String(), node.String())
	g.AddVertex(id, graph.VertexAttribute("label", nodestr))

	g.AddEdge(parent, id)

	for _, c := range node.ChildNodes() {
		offset = graphNode(g, c, id, offset)
	}

	return offset
}
