package prompt

type Prompt struct {
	user       string
	currentDir string
	ExitCode   int
}

func New() (*Prompt, error) {
	p := Prompt{}

	return &p, nil
}
