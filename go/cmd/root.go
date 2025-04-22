package cmd

import (
	"github.com/Hayao0819/nahi/cobrautils"
	"github.com/Hayao0819/zash/go/shell"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var cmdStr string
	cmd := &cobra.Command{
		Use:   "zash",
		Short: "Zash is a shell written in Go.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmdStr != "" {
				return cobrautils.CallCmd(cmd, *runCmd(), cmdStr)
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
	cmd.Flags().StringVarP(&cmdStr, "command", "c", "", "command to execute")

	cmd.AddCommand(runCmd())

	return cmd
}

func Execute() error {
	if err := rootCmd().Execute(); err != nil {
		return err
	}
	return nil
}
