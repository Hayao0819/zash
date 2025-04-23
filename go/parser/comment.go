package parser

import (
	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/lexer"
)

func (p *Parser) parseComment(cur *cursor) (*ast.Comment, error) {
	comment := &ast.Comment{
		Text: cur.next().Text,
	}

	for cur.hasNext() {
		tok := cur.peek()
		if tok.Type == lexer.TokenComment {
			comment.Text += " " + cur.next().Text
		} else {
			break
		}
	}

	return comment, nil
}
