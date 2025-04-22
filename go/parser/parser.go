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
	cmd := &ast.Command{}

	// 1. コマンド名を取得
	cmd.Name = cur.next().Text

	// 2. 引数の処理
	for cur.hasNext() {
		tok := cur.peek()

		// 空白やクォート記号はスキップ
		if tok.Type == lexer.TokenWhitespace || tok.Type == lexer.TokenQuoteChar {
			cur.next()
			continue
		}

		// 引数として使えるトークンなら処理
		if tok.Type == lexer.TokenEscapeChar || tok.Type == lexer.TokenString {
			arg := ""
			for cur.hasNext() {
				t := cur.peek()
				if t.Type != lexer.TokenEscapeChar && t.Type != lexer.TokenString {
					break
				}
				arg += cur.next().String()
			}

			if cmd.CommandSuffix == nil {
				cmd.CommandSuffix = &ast.CommandSuffix{}
			}
			cmd.CommandSuffix.Args = append(cmd.CommandSuffix.Args, arg)
		} else {
			// 予期しないトークンならスキップ
			cur.next()
		}
	}

	return cmd, nil
}
