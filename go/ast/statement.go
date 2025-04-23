package ast

type Stmt struct {
	Comments   []*Comment
	Cmd        *Command
	Redirs []*Redirection
}
