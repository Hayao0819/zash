package executer

type Executer interface {
	Exec(argv []string) (int, error)
}
