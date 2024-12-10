package main

import (
	"fmt"
	"os"
)

func handleErr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("Error: %+v\n", err)
	os.Exit(1)
}
