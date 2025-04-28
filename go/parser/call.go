package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/internal/logmgr"
	"github.com/Hayao0819/zash/go/lexer"
)

func (p *Parser) parseCommandCall(cur *cursor) (*ast.Command, error) {
	cmd := &ast.Command{
		Name:   cur.next().Text,
		Suffix: &ast.CommandSuffix{},
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

		CheckNextTokenLoop:
			for cur.hasNext() {
				t := cur.peek()

				switch t.Type {
				case lexer.TokenEscapeChar, lexer.TokenString, lexer.TokenQuotedString:
					file += cur.next().String()
				case lexer.TokenAnd:
					cur.next() // &記号の分

					if nextToAnd := cur.peek(); nextToAnd.Type == lexer.TokenNumber {
						_, err := strconv.Atoi(nextToAnd.Text)
						if err != nil {
							continue
						}

						// TODO: >&2みたいな数値でのリダイレクトに対応する
						logmgr.Parser().Warn("Redirecting to FD is not supported yet", "op", op, "file", file)

						cur.next() // >&2の数値の部分
					}
				default:
					break CheckNextTokenLoop
				}
			}
			logmgr.Parser().Info("ParserRedirectFile", "op", op, "file", file)

			// 構文エラー：ファイル名が指定されていない
			if strings.TrimSpace(file) == "" {
				return nil, fmt.Errorf("syntax error near unexpected token `%s`: missing file name", op)
			}

			cmd.Suffix.Redirections = append(cmd.Suffix.Redirections, &ast.Redirection{
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
			cmd.Suffix.Args = append(cmd.Suffix.Args, arg)

		default:
			cur.next() // skip unexpected
		}
	}

	return cmd, nil
}
