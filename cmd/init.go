package cmd

import (
	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	initCmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize a new, empty repository.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			repo.Init(path)
		},
	}
	RootCmd.AddCommand(initCmd)
}
