package cmd

import (
	"fmt"

	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	var isCreateObject bool
	tagCmd := &cobra.Command{
		Use:   "tag [name] [object]",
		Short: "List and create tags.",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			tagName := ""
			objSHA := ""
			if len(args) >= 1 {
				tagName = args[0]
			}
			if len(args) >= 2 {
				objSHA = args[1]
			}
			tags := repo.ShowOrNewTag(tagName, objSHA, isCreateObject)
			fmt.Println(tags)
		},
	}
	tagCmd.Flags().BoolVarP(&isCreateObject, "add", "a", false, "Whether to create a tag object")
	RootCmd.AddCommand(tagCmd)
}
