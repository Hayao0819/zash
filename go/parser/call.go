package parser

import (
	"fmt"
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
			redir, err := p.parseRedirection(cur)
			if err != nil {
				return nil, err
			}
			cmd.Suffix.Redirections = append(cmd.Suffix.Redirections, redir)

		case lexer.TokenEscapeChar, lexer.TokenString, lexer.TokenQuotedString:
			arg := p.parseArgument(cur)
			cmd.Suffix.Args = append(cmd.Suffix.Args, arg)

		default:
			cur.next() // skip unexpected tokens
		}
	}

	return cmd, nil
}

// parseRedirection parses a redirection (>, >>, <, etc.) and its target
func (p *Parser) parseRedirection(cur *cursor) (*ast.Redirection, error) {
	op := cur.next().Text

	// Skip whitespace and quote chars
	for cur.hasNext() && (cur.peek().Type == lexer.TokenWhitespace || cur.peek().Type == lexer.TokenQuoteChar) {
		cur.next()
	}

	file := ""

ParseFileLoop:
	for cur.hasNext() {
		t := cur.peek()

		switch t.Type {
		case lexer.TokenEscapeChar, lexer.TokenString, lexer.TokenQuotedString:
			file += cur.next().String()

		case lexer.TokenAnd:
			cur.next() // consume &
			if cur.hasNext() && cur.peek().Type == lexer.TokenNumber {
				// Example: >&2
				logmgr.Parser().Warn("Redirecting to FD is not supported yet", "op", op)
				cur.next() // consume the number
			}

		default:
			break ParseFileLoop
		}
	}

	logmgr.Parser().Info("ParserRedirectFile", "op", op, "file", file)

	if strings.TrimSpace(file) == "" {
		return nil, fmt.Errorf("syntax error near unexpected token `%s`: missing file name", op)
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
