package cmd

import (
	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	checkoutCmd := &cobra.Command{
		Use:   "checkout [commit] [dir]",
		Short: "Checkout a commit inside of a directory.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			objSHA := args[0]
			dir := args[1]
			repo.Checkout(objSHA, dir)
		},
	}
	RootCmd.AddCommand(checkoutCmd)
}
