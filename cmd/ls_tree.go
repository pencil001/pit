package cmd

import (
	"fmt"

	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	logCmd := &cobra.Command{
		Use:   "ls-tree [object]",
		Short: "Pretty-print a tree object.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			objSHA := args[0]
			tree := repo.ListTree(objSHA)
			fmt.Println(tree)
		},
	}
	RootCmd.AddCommand(logCmd)
}
