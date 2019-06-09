package cmd

import (
	"fmt"

	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	showRefCmd := &cobra.Command{
		Use:   "show-ref",
		Short: "List references.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			refs := repo.ShowRefs()
			fmt.Println(refs)
		},
	}
	RootCmd.AddCommand(showRefCmd)
}
