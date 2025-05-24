package cmd

import (
	"fmt"
	"os"

	"treereconstruction/algorithms"
	"treereconstruction/io"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(compareCmd)
}

var compareCmd = &cobra.Command{
	Use:   "compare <file1> <file2>",
	Short: "Compare two tree output files",
	Long:  `Compare two tree output files to check if they represent the same topology (structure), ignoring node names/indexes.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		content1, err := os.ReadFile(file1)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file1, err)
			return
		}

		file2 := args[1]
		content2, err := os.ReadFile(file2)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file2, err)
			return
		}

		tree1, err := io.ParseNeighborList(string(content1))
		if err != nil {
			fmt.Printf("Error parsing tree from %s: %v\n", file1, err)
			return
		}

		tree2, err := io.ParseNeighborList(string(content2))
		if err != nil {
			fmt.Printf("Error parsing tree from %s: %v\n", file2, err)
			return
		}

		err = tree1.ValidateTree()
		if err != nil {
			fmt.Printf("Error: tree from %s is invalid: %v\n", file1, err)
			return
		}

		err = tree2.ValidateTree()
		if err != nil {
			fmt.Printf("Error: tree from %s is invalid: %v\n", file2, err)
			return
		}

		if algorithms.CompareTreeTopology(tree1, tree2) {
			fmt.Printf("✓ Trees have the same topology\n")
			fmt.Printf("  %s: %d nodes, %d edges\n", file1, len(tree1.Nodes), len(tree1.AllEdges))
			fmt.Printf("  %s: %d nodes, %d edges\n", file2, len(tree2.Nodes), len(tree2.AllEdges))
		} else {
			fmt.Printf("✗ Trees have different topologies\n")
			fmt.Printf("  %s: %d nodes, %d edges\n", file1, len(tree1.Nodes), len(tree1.AllEdges))
			fmt.Printf("  %s: %d nodes, %d edges\n", file2, len(tree2.Nodes), len(tree2.AllEdges))
		}
	},
}
