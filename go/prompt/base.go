package prompt

type Prompt struct {
	user       string
	currentDir string
	exitCode   int
}

func New() (*Prompt, error) {
	p := Prompt{}
	if err := p.Update(); err != nil {
		return nil, err
	}

	return &p, nil
}
