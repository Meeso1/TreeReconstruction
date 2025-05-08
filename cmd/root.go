package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "treereconstruction",
	Short: "TreeReconstruction CLI tool",
	Long:  `A minimal CLI application for TreeReconstruction project.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is the action that will be executed when the command is called without any subcommands
		fmt.Println("Welcome to TreeReconstruction CLI!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you can define flags and configuration settings
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.treereconstruction.yaml)")
	
	// Example of a boolean flag
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
} 