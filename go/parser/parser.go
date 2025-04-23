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

func (p *Parser) Parse() (*ast.Stmt, error) {
	if len(p.tokens) == 0 {
		return nil, nil
	}

	s := &ast.Stmt{}

	cur := &cursor{tokens: p.tokens}
	first := cur.peek()

	switch first.Type {
	// case lexer.TokenIf:
	// 	return p.parseIf(cur)
	case lexer.TokenComment:
		c, err := p.parseComment(cur)
		if err != nil {
			return nil, err
		}
		if s.Comments == nil {
			s.Comments = []*ast.Comment{}
		}
		s.Comments = append(s.Comments, c)
	default:
		c, err := p.parseCommandCall(cur)
		if err != nil {
			return nil, err
		}
		s.Cmd = c
	}
	return s, nil
}
