package parser

import (
	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/lexer"
)

type Parser struct {
	tokens []lexer.Token
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (*ast.Command, error) {
	if len(p.tokens) == 0 {
		return nil, nil
	}

	cur := &cursor{tokens: p.tokens}
	first := cur.peek()

	switch first.Type {
	// case lexer.TokenIf:
	// 	return p.parseIf(cur)
	default:
		return p.parseCommandCall(cur)
	}
}
