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
		Name:   "",
		Suffix: &ast.CommandSuffix{},
	}

	for cur.hasNext() {
		tok := cur.peek()

		switch tok.Type {
		case lexer.TokenWhitespace, lexer.TokenQuoteChar:
			cur.next() // skip

		case lexer.TokenRedirection:
			logmgr.Parser().Debug("ParserFoundRedirection", "op", tok.Text)
			redir, err := p.parseRedirection(cur)
			if err != nil {
				return nil, err
			}
			cmd.Suffix.Redirections = append(cmd.Suffix.Redirections, redir)

		case lexer.TokenEscapeChar, lexer.TokenString, lexer.TokenQuotedString:
			if cmd.Name == "" {
				cmd.Name = cur.next().String()
			} else {
				arg := p.parseArgument(cur)
				cmd.Suffix.Args = append(cmd.Suffix.Args, arg)
			}
		case lexer.TokenPipe:
			cur.next() // consume pipe
			if !cur.hasNext() {
				return nil, fmt.Errorf("syntax error: expected command after pipe")
			}
			nextCmd, err := p.parseCommandCall(cur)
			if err != nil {
				return nil, err
			}
			cmd.Next = nextCmd
			logmgr.Parser().Debug("ParserFoundPipe", "nextCmd", nextCmd.Name)
		default:
			cur.next() // skip unexpected tokens
		}
	}

	return cmd, nil
}

// parseRedirection parses a redirection (>, >>, <, etc.) and its target
func (p *Parser) parseRedirection(cur *cursor) (*ast.Redirection, error) {
	op := cur.next().Text // consume redirection operator (e.g., >, <)

	// Skip whitespace and quote characters
	for cur.hasNext() {
		switch cur.peek().Type {
		case lexer.TokenWhitespace, lexer.TokenQuoteChar:
			cur.next()
		default:
			goto FindFile
		}
	}

FindFile:
	// ファイル名・またはファイルディスクリプタを取得する
	file := ""

	if !cur.hasNext() {
		return nil, fmt.Errorf("syntax error: expected filename after '%s'", op)
	}

ReadFileLoop:
	for cur.hasNext() {
		tok := cur.peek()

		switch tok.Type {
		case lexer.TokenEscapeChar, lexer.TokenString, lexer.TokenQuotedString:
			file += cur.next().String()
		case lexer.TokenAnd:
			cur.next() // consume '&'
			if cur.hasNext() && cur.peek().Type == lexer.TokenNumber {
				numTok := cur.next()
				// 数値じゃなかったらエラー扱い
				if _, err := strconv.Atoi(numTok.Text); err != nil {
					return nil, fmt.Errorf("syntax error: invalid file descriptor after '>&'")
				}
				logmgr.Parser().Warn("Redirecting to FD is not supported yet", "op", op, "fd", numTok.Text)
			} else {
				return nil, fmt.Errorf("syntax error: expected file descriptor after '>&'")
			}
		case lexer.TokenRedirection:
			// 連続リダイレクト禁止
			return nil, fmt.Errorf("syntax error: unexpected '%s' after redirection '%s'", tok.Text, op)
		default:
			break ReadFileLoop
		}
	}

	if strings.TrimSpace(file) == "" {
		return nil, fmt.Errorf("syntax error: missing filename after '%s'", op)
	}

	return &ast.Redirection{
		Operator: op,
		File:     file,
	}, nil
}

// parseArgument parses a single argument (could be escaped strings, quoted strings, etc.)
func (p *Parser) parseArgument(cur *cursor) string {
	arg := ""

	for cur.hasNext() {
		t := cur.peek()
		if t.Type != lexer.TokenEscapeChar && t.Type != lexer.TokenString && t.Type != lexer.TokenQuotedString {
			break
		}
		arg += cur.next().String()
	}

	return arg
}
