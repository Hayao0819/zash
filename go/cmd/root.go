package cmd

import (
	"github.com/Hayao0819/nahi/cobrautils"
	"github.com/Hayao0819/zash/go/internal/logmgr"
	"github.com/Hayao0819/zash/go/shell"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var optCmdStr string
	var optDebug bool
	cmd := &cobra.Command{
		Use:   "zash",
		Short: "Zash is a shell written in Go.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if optDebug {
				logmgr.EnableAll()
				logmgr.Shell().Info("debug mode enabled")
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if optCmdStr != "" {
				return cobrautils.CallCmd(cmd, *runCmd(), optCmdStr)
			}

			s, err := shell.New()
			if err != nil {
				return err
			}
			
			s.StartInteractive()
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	// cmd.Flags().StringP("command", "c", "", "command to execute")
	cmd.Flags().StringVarP(&optCmdStr, "command", "c", "", "command to execute")
	cmd.Flags().BoolVarP(&optDebug, "debug", "d", false, "debug mode")

	cmd.AddCommand(runCmd())

	return cmd
}
