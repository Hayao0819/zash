package parser

import "github.com/Hayao0819/zash/go/lexer"

type cursor struct {
	processed int
	tokens    []lexer.Token
}

func (c *cursor) next() lexer.Token {
	if c.processed >= len(c.tokens) {
		return lexer.Token{}
	}
	token := c.tokens[c.processed]
	c.processed++
	return token
}

func (c *cursor) peek() lexer.Token {
	if c.processed >= len(c.tokens) {
		return lexer.Token{}
	}
	return c.tokens[c.processed]
}
func (c *cursor) left() []lexer.Token {
	if c.processed >= len(c.tokens) {
		return []lexer.Token{}
	}
	return c.tokens[c.processed:]
}
func (c *cursor) hasNext() bool {
	return c.processed < len(c.tokens)
}
