package ast

import "encoding/json"

type Command struct {
	Name   string
	Suffix *CommandSuffix
	Next   *Command
}

func (c *Command) Argv() []string {
	argv := []string{c.Name}
	if c.Suffix != nil {
		argv = append(argv, c.Suffix.Args...)
	}
	return argv
}

type CommandSuffix struct {
	Args         []string
	Redirections []*Redirection
}

func (c *Command) JSON() ([]byte, error) {
	return json.Marshal(c)
}

type Redirection struct {
	Operator string // ">", ">>", "<" など
	File     string // リダイレクト先（ファイル名）
}
