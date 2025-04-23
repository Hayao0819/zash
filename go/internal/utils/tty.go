package utils

import (
	"fmt"

	"github.com/mattn/go-tty"
)

func TTYPrintln(t *tty.TTY, a ...any) {
	fmt.Fprintln(t.Output(), a...)
}

func TTYPrint(t *tty.TTY, a ...any) {
	fmt.Fprintf(t.Output(), "%s", a...)
}
func TTYPrintf(t *tty.TTY, format string, a ...any) {
	fmt.Fprintf(t.Output(), format, a...)
}
