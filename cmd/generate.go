package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"treereconstruction/algorithms"
	"treereconstruction/io"

	"github.com/spf13/cobra"
)

var (
	numLeaves             int
	outputPrefix          string
	seed                  int64
	chainExtensionProb    float64
	connectToExistingProb float64
)

func init() {
	generateCmd.Flags().IntVarP(&numLeaves, "leaves", "l", 4, "Number of leaves in the generated tree")
	generateCmd.Flags().StringVarP(&outputPrefix, "prefix", "p", "", "Output file prefix (required)")
	generateCmd.Flags().Int64VarP(&seed, "seed", "s", 0, "Random seed (0 for current time)")
	generateCmd.Flags().Float64VarP(&chainExtensionProb, "chain-prob", "c", 0, "Probability of extending leaf chains (0 = no chains, 0.1 = 10% chance per extension)")
	generateCmd.Flags().Float64VarP(&connectToExistingProb, "connect-prob", "x", 0.25, "Probability of connecting new leaves to existing nodes instead of splitting edges (default: 0.25)")
	generateCmd.MarkFlagRequired("prefix")

	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate example inputs and outputs",
	Long:  `Generate a random tree with specified number of leaves and create corresponding distance matrix input and tree output files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if numLeaves < 2 {
			fmt.Printf("Number of leaves must be at least 2, got %d\n", numLeaves)
			return
		}

		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		//fmt.Printf("Using random seed: %d\n", seed)

		inputFile := fmt.Sprintf("%s-%d.input.txt", outputPrefix, numLeaves)
		if _, err := os.Stat(inputFile); err == nil {
			fmt.Printf("Error: input file %s already exists\n", inputFile)
			return
		}

		outputFile := fmt.Sprintf("%s-%d.output.txt", outputPrefix, numLeaves)
		if _, err := os.Stat(outputFile); err == nil {
			fmt.Printf("Error: output file %s already exists\n", outputFile)
			return
		}

		//fmt.Printf("Generating random tree with %d leaves...\n", numLeaves)
		tree, err := algorithms.GenerateRandomTree(numLeaves, seed, chainExtensionProb, connectToExistingProb)
		if err != nil {
			fmt.Printf("Error generating tree: %v\n", err)
			return
		}

		err = tree.ValidateTree()
		if err != nil {
			fmt.Printf("Error: generated tree is invalid: %v\n", err)
			return
		}

		actualLeaves := len(algorithms.GetLeafNodes(tree))
		if actualLeaves != numLeaves {
			fmt.Printf("Error: generated tree has %d leaves, expected %d\n", actualLeaves, numLeaves)
			return
		}

		//fmt.Printf("Calculating distance matrix...\n")
		distanceMatrix, err := algorithms.CalculateDistanceMatrix(tree)
		if err != nil {
			fmt.Printf("Error calculating distance matrix: %v\n", err)
			return
		}

		matrixCSV := algorithms.FormatDistanceMatrix(distanceMatrix)

		//fmt.Printf("Serializing tree...\n")
		serializedTree, err := io.SerializeGraph(tree, io.SerializationTypeNeighborLists)
		if err != nil {
			fmt.Printf("Error serializing tree: %v\n", err)
			return
		}

		inputDir := filepath.Dir(inputFile)
		if inputDir != "." {
			if err := os.MkdirAll(inputDir, 0755); err != nil {
				fmt.Printf("Error creating input directory: %v\n", err)
				return
			}
		}

		outputDir := filepath.Dir(outputFile)
		if outputDir != "." {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Printf("Error creating output directory: %v\n", err)
				return
			}
		}

		//fmt.Printf("Writing distance matrix to %s...\n", inputFile)
		err = os.WriteFile(inputFile, []byte(matrixCSV), 0644)
		if err != nil {
			fmt.Printf("Error writing input file: %v\n", err)
			return
		}

		//fmt.Printf("Writing tree to %s...\n", outputFile)
		err = os.WriteFile(outputFile, []byte(serializedTree), 0644)
		if err != nil {
			fmt.Printf("Error writing output file: %v\n", err)
			return
		}

		fmt.Printf("Successfully generated:\n")
		fmt.Printf("  Input file:  %s (%dx%d distance matrix)\n", inputFile, numLeaves, numLeaves)
		fmt.Printf("  Output file: %s (neighbor lists format)\n", outputFile)
		fmt.Printf("  Random seed: %d\n", seed)

		fmt.Printf("  Tree summary:\n")
		treeSummary := io.GetTreeSummary(tree)
		for _, line := range treeSummary {
			fmt.Printf("    %s\n", line)
		}
	},
}
