package main

import (
	"fmt"

	"github.com/Hayao0819/zash/go/internal/utils"
	"github.com/Hayao0819/zash/go/lexer"
	"github.com/Hayao0819/zash/go/shell"
)

func main() {
	fmt.Println("Welcome to Zash!")

	tokens, _ := lexer.NewLexer("echo \"hello world\"").ReadAll()
	for _, token := range tokens {
		fmt.Println("t: (", token,")")
	}

	_, err := shell.New()
	utils.HandleErr(err)
	// s.StartInteractive()
}
