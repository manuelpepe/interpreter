package ast

import (
	"bytes"
	"strings"

	"github.com/manuelpepe/interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	ChildNodes() []Node
}

type Statement interface {
	Node
	statementNode() // helps go type-checker
}

type Expression interface {
	Node
	expressionNode() // helps go type-checker
}

// Program is always the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) ChildNodes() []Node {
	nodes := make([]Node, len(p.Statements))
	for ix := range p.Statements {
		nodes[ix] = p.Statements[ix]
	}
	return nodes
}

// LetStatement is used to define variables
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) ChildNodes() []Node {
	return []Node{s.Name, s.Value}
}
func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")
	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Identifiers are things like variable and function names, and constants like numbers
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) ChildNodes() []Node {
	return []Node{}
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (s *ReturnStatement) ChildNodes() []Node {
	return []Node{s.ReturnValue}
}
func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	out.WriteString(" ")
	if s.ReturnValue != nil {
		out.WriteString(s.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement allows lines that are only expressions to be added to the program
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (s *ExpressionStatement) ChildNodes() []Node {
	return []Node{s.Expression}
}
func (s *ExpressionStatement) statementNode()       {}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ChildNodes() []Node {
	return []Node{}
}
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // prefix token, eg. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) ChildNodes() []Node {
	return []Node{pe.Right}
}
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // infix token, eg. + / -
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) ChildNodes() []Node {
	return []Node{ie.Left, ie.Right}
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) ChildNodes() []Node {
	return []Node{}
}
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// BlockStatements are a series of Statements enclosed between curly braces { s1; s2; }
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) ChildNodes() []Node {
	nodes := make([]Node, len(bs.Statements))
	for ix := range bs.Statements {
		nodes[ix] = bs.Statements[ix]
	}
	return nodes
}
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) ChildNodes() []Node {
	nodes := make([]Node, 0)
	nodes = append(nodes, ie.Condition)
	nodes = append(nodes, ie.Consequence)
	if ie.Alternative != nil {
		nodes = append(nodes, ie.Alternative)
	}
	return nodes
}
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) ChildNodes() []Node {
	nodes := make([]Node, len(fl.Parameters)+1)
	for ix := range fl.Parameters {
		nodes[ix] = fl.Parameters[ix]
	}
	nodes[len(fl.Parameters)] = fl.Body
	return nodes
}
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := make([]string, len(fl.Parameters))
	for ix, p := range fl.Parameters {
		params[ix] = p.String()
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // the ( token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (c *CallExpression) ChildNodes() []Node {
	nodes := make([]Node, len(c.Arguments)+1)
	nodes[0] = c.Function
	for ix := range c.Arguments {
		nodes[ix+1] = c.Arguments[ix]
	}
	return nodes
}
func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var out bytes.Buffer

	args := make([]string, len(c.Arguments))
	for ix, a := range c.Arguments {
		args[ix] = a.String()
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
