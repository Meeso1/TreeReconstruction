package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

func init() {
	processCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path (required)")
	processCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	processCmd.MarkFlagRequired("input")
	
	rootCmd.AddCommand(processCmd)
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process a tree file",
	Long:  `Process a tree file with the specified options.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Processing file: %s\n", inputFile)
		if outputFile != "" {
			fmt.Printf("Output file: %s\n", outputFile)
		} else {
			fmt.Println("Using standard output")
		}
		
		// This is where you would put the logic to process the file
		fmt.Println("Processing complete!")
	},
} 