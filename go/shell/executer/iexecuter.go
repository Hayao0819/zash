package executer

type Executer interface {
	Exec(argv []string) error
}
