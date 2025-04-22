package ast

import "encoding/json"

type Command struct {
	Name          string
	CommandSuffix *CommandSuffix
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
