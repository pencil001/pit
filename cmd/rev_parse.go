package cmd

import (
	"fmt"
	"log"

	"github.com/pencil001/pit/repo"
	"github.com/pencil001/pit/util"
	"github.com/spf13/cobra"
)

func init() {
	var revType string
	revParseCmd := &cobra.Command{
		Use:   "rev-parse [name]",
		Short: "Parse revision (or other objects) identifiers.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			objRev := args[0]
			if revType != "" && !util.ObjectIn([]string{repo.TypeBlob, repo.TypeCommit, repo.TypeTag, repo.TypeTree}, revType) {
				log.Panicf("Unknown type: %v", revType)
			}
			hash := repo.RevParse(objRev, revType)
			fmt.Println(hash)
		},
	}
	revParseCmd.Flags().StringVarP(&revType, "type", "t", "", "Specify the expected type")
	RootCmd.AddCommand(revParseCmd)
}
