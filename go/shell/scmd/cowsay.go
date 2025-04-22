package scmd

import cowsay "github.com/Code-Hex/Neo-cowsay/v2"

var cowsayCmd = InternalCmd{
	Name: "cowsay",
	Func: func(e Executer, args []string) Result {
		if len(args) == 0 {
			return Result{
				err:      nil,
				exitcode: 0,
			}
		}

		s, err := cowsay.Say(args[0], cowsay.Type("default"))
		if err != nil {
			return Result{
				err:      err,
				exitcode: 1,
			}
		}

		e.Puts(s)
		return Result{
			err:      nil,
			exitcode: 0,
		}
	},
}

func init() {
	internalCmds = append(internalCmds, cowsayCmd)
}
