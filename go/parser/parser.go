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

	processed := 0
	left := func() []lexer.Token {
		// return p.tokens[processed:]
		if processed >= len(p.tokens) {
			return []lexer.Token{}
		}
		return p.tokens[processed:]
	}
	next := func() lexer.Token {
		if len(left()) == 0 {
			return lexer.Token{}
		}
		token := left()[0]
		processed++
		return token
	}
	// peek は次のトークンを返すが、processedは進めない
	peek := func() lexer.Token {
		if len(left()) == 0 {
			return lexer.Token{}
		}
		// fmt.Println("peek", left()[0].Text, "processed", processed)
		return left()[0]
	}

	// 1. コマンド名を取得
	cmd := &ast.Command{}
	cmd.Name = next().Text

	// 2. コマンドの引数を取得
	if len(left()) > 0 {
		cmd.CommandSuffix = &ast.CommandSuffix{}
		currentText := ""
		for _, token := range left() {
			if token.Type == lexer.TokenWhitespace || token.Type == lexer.TokenQuoteChar {
				next()
				continue
			}
			if token.Type == lexer.TokenEscapeChar || token.Type == lexer.TokenString {
				currentText += token.String()
				processed++
				if peek().Type == lexer.TokenWhitespace || peek().Type == lexer.TokenEOT {
					// fmt.Println("append currentText", currentText)
					cmd.CommandSuffix.Args = append(cmd.CommandSuffix.Args, currentText)
					currentText = ""
				}
			}
			continue
		}
	}

	return cmd, nil
}
