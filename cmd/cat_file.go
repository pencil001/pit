package cmd

import (
	"fmt"
	"log"

	"github.com/pencil001/pit/repo"
	"github.com/pencil001/pit/util"
	"github.com/spf13/cobra"
)

func init() {
	catFileCmd := &cobra.Command{
		Use:   "cat-file [type] [object]",
		Short: "Provide content of repository objects",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			objType := args[0]
			objSHA := args[1]
			if !util.ObjectIn([]string{repo.TypeBlob, repo.TypeCommit, repo.TypeTag, repo.TypeTree}, objType) {
				log.Panicf("Unknown type: %v", objType)
			}
			content := repo.Cat(objType, objSHA)
			fmt.Println(content)
		},
	}
	RootCmd.AddCommand(catFileCmd)
}
