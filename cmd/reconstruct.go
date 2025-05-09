package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"treereconstruction/io"
	"os"
)

var (
	inputFile  string
	outputFile string
)

func init() {
	reconstructCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path (required)")
	reconstructCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	reconstructCmd.MarkFlagRequired("input")
	
	rootCmd.AddCommand(reconstructCmd)
}

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct",
	Short: "Reconstruct a tree",
	Long:  `Reconstruct a tree from distance matrix`,
	Run: func(cmd *cobra.Command, args []string) {
		fileContent, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		
		matrix, err := io.ParseMatrix(string(fileContent))
		if err != nil {
			fmt.Printf("Error parsing matrix: %v\n", err)
			return
		}
		
		fmt.Printf("Matrix: %v\n", matrix)
	},
} 