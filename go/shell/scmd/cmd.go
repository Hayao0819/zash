package scmd

type InternalCmd struct {
	Name string
	Func func(e Executer, args []string) Result
	// State map[string]any
}
