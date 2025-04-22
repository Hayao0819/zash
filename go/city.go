package main

import (
	"github.com/Hayao0819/zash/go/internal/utils"
	"github.com/Hayao0819/zash/go/shell"
)

func main() {
	s, err := shell.New()
	utils.HandleErr(err)
	s.StartInteractive()
}
