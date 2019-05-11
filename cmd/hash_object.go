package cmd

import (
	"fmt"

	"github.com/pencil001/pit/util"

	"github.com/pencil001/pit/repo"
	"github.com/spf13/cobra"
)

func init() {
	var isStore bool
	var objType string
	hashObjectCmd := &cobra.Command{
		Use:   "hash-object [file]",
		Short: "Compute object ID and optionally creates a blob from a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !util.ObjectIn([]string{repo.TypeBlob, repo.TypeCommit, repo.TypeTag, repo.TypeTree}, objType) {
				fmt.Errorf("Unknown type: %v", objType)
			}
			file := args[0]
			hash := repo.Hash(file, objType, isStore)
			fmt.Println(hash)
		},
	}
	hashObjectCmd.Flags().BoolVarP(&isStore, "write", "w", false, "Actually write the object into the database")
	hashObjectCmd.Flags().StringVarP(&objType, "type", "t", repo.TypeBlob, "Specify the type")
	RootCmd.AddCommand(hashObjectCmd)
}
