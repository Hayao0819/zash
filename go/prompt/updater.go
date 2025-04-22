package prompt

import (
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

func (p *Prompt) Update() error {
	eg := new(errgroup.Group)
	eg.Go(p.UpdateUser)
	eg.Go(p.UpdateCurrentDir)
	return eg.Wait()
}

func (p *Prompt) UpdateUser() error {
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("user not found")
	}
	p.user = user
	return nil
}

func (p *Prompt) UpdateCurrentDir() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	p.currentDir = currentDir
	return nil
}

func (p *Prompt) SetExitCode(code int) {
	p.exitCode = code
}
