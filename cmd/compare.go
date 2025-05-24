package cmd

import (
	"fmt"
	"os"

	"treereconstruction/algorithms"
	"treereconstruction/io"

	"github.com/spf13/cobra"
)

type CompareResult struct {
	TopologiesMatch bool
	Tree1Summary    []string
	Tree2Summary    []string
	Error           error
}

func init() {
	rootCmd.AddCommand(compareCmd)
}

func runCompareCommand(file1, file2 string) CompareResult {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return CompareResult{Error: fmt.Errorf("error reading file %s: %v", file1, err)}
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return CompareResult{Error: fmt.Errorf("error reading file %s: %v", file2, err)}
	}

	tree1, err := io.ParseNeighborList(string(content1))
	if err != nil {
		return CompareResult{Error: fmt.Errorf("error parsing tree from %s: %v", file1, err)}
	}

	tree2, err := io.ParseNeighborList(string(content2))
	if err != nil {
		return CompareResult{Error: fmt.Errorf("error parsing tree from %s: %v", file2, err)}
	}

	err = tree1.ValidateTree()
	if err != nil {
		return CompareResult{Error: fmt.Errorf("tree from %s is invalid: %v", file1, err)}
	}

	err = tree2.ValidateTree()
	if err != nil {
		return CompareResult{Error: fmt.Errorf("tree from %s is invalid: %v", file2, err)}
	}

	topologiesMatch := algorithms.CompareTreeTopology(tree1, tree2)
	tree1Summary := io.GetTreeSummary(tree1)
	tree2Summary := io.GetTreeSummary(tree2)

	return CompareResult{
		TopologiesMatch: topologiesMatch,
		Tree1Summary:    tree1Summary,
		Tree2Summary:    tree2Summary,
		Error:           nil,
	}
}

var compareCmd = &cobra.Command{
	Use:   "compare <file1> <file2>",
	Short: "Compare two tree output files",
	Long:  `Compare two tree output files to check if they represent the same topology (structure), ignoring node names/indexes.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		file2 := args[1]

		result := runCompareCommand(file1, file2)
		if result.Error != nil {
			fmt.Printf("%v\n", result.Error)
			return
		}

		if result.TopologiesMatch {
			fmt.Printf("✓ Trees have the same topology\n")
			for _, line := range result.Tree1Summary {
				fmt.Printf("  %s\n", line)
			}
		} else {
			fmt.Printf("✗ Trees have different topologies\n")

			fmt.Printf("  %s:\n", file1)
			for _, line := range result.Tree1Summary {
				fmt.Printf("    %s\n", line)
			}

			fmt.Printf("  %s:\n", file2)
			for _, line := range result.Tree2Summary {
				fmt.Printf("    %s\n", line)
			}
		}
	},
}
