package parser

import (
	"fmt"
	"strings"

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
			// 構文エラー：ファイル名が指定されていない
			if strings.TrimSpace(file) == "" {
				return nil, fmt.Errorf("syntax error near unexpected token `%s`: missing file name", op)
			}

			cmd.CommandSuffix.Redirections = append(cmd.CommandSuffix.Redirections, &ast.Redirection{
				Operator: op,
				File:     file,
			})

		case lexer.TokenEscapeChar, lexer.TokenString, lexer.TokenQuotedString:
			arg := ""
			for cur.hasNext() {
				t := cur.peek()
				if t.Type != lexer.TokenEscapeChar && t.Type != lexer.TokenString && t.Type != lexer.TokenQuotedString {
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
