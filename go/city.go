package main

import (
	"fmt"

	"github.com/Hayao0819/zash/go/internal/utils"
	"github.com/Hayao0819/zash/go/shell"
)

func main() {
	fmt.Println("Welcome to Zash!")

	s, err := shell.New()
	utils.HandleErr(err)
	s.StartInteractive()
}
