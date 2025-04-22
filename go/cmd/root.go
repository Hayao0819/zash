package cmd

import (
	"github.com/Hayao0819/zash/go/internal/utils"
	"github.com/Hayao0819/zash/go/shell"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zash",
		Short: "Zash is a shell written in Go.",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := shell.New()
			utils.HandleErr(err)
			s.StartInteractive()
		},
	}

	return cmd
}

func Execute() error {
	if err := rootCmd().Execute(); err != nil {
		return err
	}
	return nil
}
