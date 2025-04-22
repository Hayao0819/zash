package cmd

import (
	"github.com/Hayao0819/zash/go/shell"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zash",
		Short: "Zash is a shell written in Go.",
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.AddCommand(runCmd())

	return cmd
}

func Execute() error {
	if err := rootCmd().Execute(); err != nil {
		return err
	}
	return nil
}
