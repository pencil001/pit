package cmd

import (
	"fmt"

	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	logCmd := &cobra.Command{
		Use:   "log [commit]",
		Short: "Display history of a given commit.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			head := "HEAD"
			if len(args) == 1 {
				head = args[0]
			}
			log := repo.Log(head)
			fmt.Println(log)
		},
	}
	RootCmd.AddCommand(logCmd)
}
