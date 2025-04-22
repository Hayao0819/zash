package parser

import (
	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/lexer"
)

func (p *Parser) parseCommandCall(cur *cursor) (*ast.Command, error) {
	cmd := &ast.Command{
		Name:          cur.next().Text,
		CommandSuffix: &ast.CommandSuffix{},
	}

	for cur.hasNext() {
		tok := cur.peek()

		switch tok.Type {
		case lexer.TokenWhitespace, lexer.TokenQuoteChar:
			cur.next() // skip
		case lexer.TokenRedirection:
			// Handle redirection
			op := cur.next().Text
			for cur.hasNext() && (cur.peek().Type == lexer.TokenWhitespace || cur.peek().Type == lexer.TokenQuoteChar) {
				cur.next()
			}

			file := ""
			for cur.hasNext() {
				t := cur.peek()
				if t.Type != lexer.TokenEscapeChar && t.Type != lexer.TokenString {
					break
				}
				file += cur.next().String()
			}

			cmd.CommandSuffix.Redirections = append(cmd.CommandSuffix.Redirections, &ast.Redirection{
				Operator: op,
				File:     file,
			})

		case lexer.TokenEscapeChar, lexer.TokenString:
			// Handle arguments
			arg := ""
			for cur.hasNext() {
				t := cur.peek()
				if t.Type != lexer.TokenEscapeChar && t.Type != lexer.TokenString {
					break
				}
				arg += cur.next().String()
			}
			cmd.CommandSuffix.Args = append(cmd.CommandSuffix.Args, arg)

		default:
			cur.next() // skip unexpected
		}
	}

	return cmd, nil
}
