package scmd

// var typeCmd = InternalCmd{
// 	Name: "type",
// 	Func: func(t *tty.TTY) func(args []string) error {
// 		return func(args []string) error {
// 			if len(args) == 0 {
// 				return nil
// 			}
// 			for _, arg := range args {
// 				cmd := NewInternalCmds(t)
// 				if cmd.Get(arg) != nil {
// 					// t.Write([]byte("internal command\n"))
// 					fmt.Fprintln(t.Output(), "internal command")
// 					continue
// 				}

// 				p, err := exec.LookPath(arg)
// 				if err == nil {
// 					fmt.Fprintf(t.Output(), "external command: %s\n", p)
// 					continue
// 				}
// 				return fmt.Errorf("type: %s: not found", arg)
// 			}
// 			return nil
// 		}
// 	},
// }
