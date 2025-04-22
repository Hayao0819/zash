package shell

import (
	"fmt"
)

func (s *Shell) Print(str string) {
	fmt.Fprintf(s.TTY.Output(), "%s", str)
}
func (s *Shell) Println(str string) {
	fmt.Fprintln(s.TTY.Output(), str)
}

func (s *Shell) Close() {
	s.TTY.Close()
}
