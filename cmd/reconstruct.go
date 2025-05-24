package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"treereconstruction/algorithms"
	"treereconstruction/io"

	"github.com/spf13/cobra"
)

var (
	inputFile               string
	outputFile              string
	serializationTypeString string
)

type ReconstructResult struct {
	SerializedTree string
	Error          error
}

func init() {
	reconstructCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path (required)")
	reconstructCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	reconstructCmd.Flags().StringVarP(&serializationTypeString, "serialization", "s", "neighbor-lists", "Serialization type (brackets, brackets-shortened, neighbor-lists)")
	reconstructCmd.MarkFlagRequired("input")

	rootCmd.AddCommand(reconstructCmd)
}

func runReconstructCommand(inputFilePath, outputFilePath string, serializationType io.SerializationType) ReconstructResult {
	fileContent, err := os.ReadFile(inputFilePath)
	if err != nil {
		return ReconstructResult{Error: fmt.Errorf("error reading file: %v", err)}
	}

	matrix, err := io.ParseMatrix(string(fileContent))
	if err != nil {
		return ReconstructResult{Error: fmt.Errorf("error parsing matrix: %v", err)}
	}

	epsilon := 1e-10
	tree, err := algorithms.ReconstructIntTree(matrix, epsilon)
	if err != nil {
		return ReconstructResult{Error: fmt.Errorf("error reconstructing tree: %v", err)}
	}

	if !tree.IsIntegerWeighted(epsilon) {
		return ReconstructResult{Error: fmt.Errorf("tree is not integer weighted")}
	}

	serialized, err := io.SerializeGraph(tree, serializationType)
	if err != nil {
		return ReconstructResult{Error: fmt.Errorf("error serializing tree: %v", err)}
	}

	if outputFilePath != "" {
		if _, err := os.Stat(outputFilePath); err == nil {
			os.Remove(outputFilePath)
		}

		if err := os.MkdirAll(filepath.Dir(outputFilePath), 0755); err != nil {
			return ReconstructResult{Error: fmt.Errorf("error creating output directory: %v", err)}
		}

		err = os.WriteFile(outputFilePath, []byte(serialized), 0644)
		if err != nil {
			return ReconstructResult{Error: fmt.Errorf("error writing output file: %v", err)}
		}
	}

	return ReconstructResult{SerializedTree: serialized, Error: nil}
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

		result := runReconstructCommand(inputFile, outputFile, serializationType)
		if result.Error != nil {
			fmt.Printf("%v\n", result.Error)
			return
		}

		if serializationType == io.SerializationTypeNeighborLists {
			fmt.Printf("Tree:\n%v\n", result.SerializedTree)
		} else {
			fmt.Printf("Tree: %v\n", result.SerializedTree)
		}
	},
}
