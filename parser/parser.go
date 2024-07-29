package parser

import (
	"fmt"

	"github.com/manuelpepe/interpreter/ast"
	"github.com/manuelpepe/interpreter/lexer"
	"github.com/manuelpepe/interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: make([]string, 0)}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{Statements: make([]ast.Statement, 0)}
	for !p.curTokenIs(token.EOF) {
		stm := p.parseStatement()
		if stm != nil {
			prog.Statements = append(prog.Statements, stm)
		}
		p.nextToken()
	}
	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	tkn := p.curToken

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ident := p.parseIdentifier()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	value := p.parseExpression()

	return &ast.LetStatement{
		Token: tkn,
		Name:  ident,
		Value: value,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	tkn := p.curToken
	p.nextToken() // prepare for parsing expr
	value := p.parseExpression()
	return &ast.ReturnStatement{
		Token:       tkn,
		ReturnValue: value,
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek advances the parser position only if the next token is of type `t`, otherwise registers a peekError
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	err := fmt.Sprintf("expected next token to be %s, got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}
