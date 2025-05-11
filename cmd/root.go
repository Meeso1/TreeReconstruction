package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tree-reconstruction",
	Short: "Tree reconstruction CLI tool",
	Long:  `CLI application that reconstructs a tree from a distance matrix.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run `tree-reconstruction reconstruct` to reconstruct a tree, or `tree-reconstruction help` for a list of available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
