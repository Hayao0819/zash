package main

import (
	"github.com/Hayao0819/zash/go/cmd"
	"github.com/Hayao0819/zash/go/internal/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.HandleErr(err)
	}
}
