package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "",
		Args:  cobra.NoArgs,
	}
	RootCmd.AddCommand(addCmd)
}
