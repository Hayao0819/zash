package ast

import "encoding/json"

type Command struct {
	Name          string
	CommandSuffix *CommandSuffix
}

type CommandSuffix struct {
	Args []string
}

func (c *Command) JSON() ([]byte, error) {
	return json.Marshal(c)
}
