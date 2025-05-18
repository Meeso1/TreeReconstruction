package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"treereconstruction/io"
	"treereconstruction/algorithms"
	"os"
	"path/filepath"
)

var (
	inputFile  string
	outputFile string
	serializationTypeString string
)

func init() {
	reconstructCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path (required)")
	reconstructCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	reconstructCmd.Flags().StringVarP(&serializationTypeString, "serialization", "s", "neighbor-lists", "Serialization type (brackets, brackets-shortened, neighbor-lists)")
	reconstructCmd.MarkFlagRequired("input")
	
	rootCmd.AddCommand(reconstructCmd)
}

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct",
	Short: "Reconstruct a tree",
	Long:  `Reconstruct a tree from distance matrix`,
	Run: func(cmd *cobra.Command, args []string) {
		var serializationType io.SerializationType
		switch serializationTypeString {
		case "brackets":
			serializationType = io.SerializationTypeBrackets
		case "brackets-shortened":
			serializationType = io.SerializationTypeBracketsShortened
		case "neighbor-lists":
			serializationType = io.SerializationTypeNeighborLists
		default:
			fmt.Printf("Invalid serialization type: %s\n", serializationTypeString)
			return
		}

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
		
		epsilon := 1e-10
		tree, err := algorithms.ReconstructIntTree(matrix, epsilon)
		if err != nil {
			fmt.Printf("Error reconstructing tree: %v\n", err)
			return
		}

		if !tree.IsIntegerWeighted(epsilon) {
			fmt.Printf("Tree is not integer weighted\n")
			return
		}

		serialized, err := io.SerializeGraph(tree, serializationType)
		if err != nil {
			fmt.Printf("Error serializing tree: %v\n", err)
			return
		}
		
		if outputFile != "" {
			if _, err := os.Stat(outputFile); err == nil {
				os.Remove(outputFile)
			}

			if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
				fmt.Printf("Error creating output directory: %v\n", err)
				return
			}

			err = os.WriteFile(outputFile, []byte(serialized), 0644)
			if err != nil {
				fmt.Printf("Error writing output file: %v\n", err)
				return
			}
		}

		if serializationType == io.SerializationTypeNeighborLists {
			fmt.Printf("Tree:\n%v\n", serialized)
		} else {
			fmt.Printf("Tree: %v\n", serialized)
		}
	},
} 