package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the root of cmd
var RootCmd = &cobra.Command{
	Use:   "pit",
	Short: "A self-implemented git",
}

// Execute is the entrance of the cobra cmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
