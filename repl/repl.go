package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/manuelpepe/interpreter/eval"
	"github.com/manuelpepe/interpreter/lexer"
	"github.com/manuelpepe/interpreter/object"
	"github.com/manuelpepe/interpreter/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.New(l)

		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		res := eval.Eval(prog, env)
		if res != nil {
			io.WriteString(out, res.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
