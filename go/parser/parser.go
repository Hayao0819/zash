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

	var cmd ast.Command
	i := 0 // トークンの現在位置

	// トークンの取得ヘルパー
	next := func() lexer.Token {
		if i >= len(p.tokens) {
			return lexer.Token{}
		}
		tok := p.tokens[i]
		i++
		return tok
	}

	// 1. コマンド名を取得
	cmd.Name = next().Text

	// 2. 引数を取得
	for i < len(p.tokens) {
		tok := p.tokens[i]

		// 空白とクォート文字はスキップ
		if tok.Type == lexer.TokenWhitespace || tok.Type == lexer.TokenQuoteChar {
			i++
			continue
		}

		// エスケープ文字や通常の文字列を処理
		if tok.Type == lexer.TokenEscapeChar || tok.Type == lexer.TokenString {
			arg := ""

			// 連続する文字列トークンを結合
			for i < len(p.tokens) {
				t := p.tokens[i]
				if t.Type != lexer.TokenEscapeChar && t.Type != lexer.TokenString {
					break
				}
				arg += t.String()
				i++
			}

			// 次のトークンが空白またはEOTなら引数を確定
			if cmd.CommandSuffix == nil {
				cmd.CommandSuffix = &ast.CommandSuffix{}
			}
			cmd.CommandSuffix.Args = append(cmd.CommandSuffix.Args, arg)
		} else {
			i++
		}
	}

	return &cmd, nil
}
