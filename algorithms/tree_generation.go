package algorithms

import (
	"fmt"
	"math/rand"
)

// Creates a random tree with the specified number of leaves.
// The function starts with a simple path and randomly adds new leaves until
// the desired number is reached.
func GenerateRandomTree(numLeaves int, seed int64) (*Graph, error) {
	if numLeaves < 2 {
		return nil, fmt.Errorf("number of leaves must be at least 2, got %d", numLeaves)
	}

	rng := rand.New(rand.NewSource(seed))

	graph := &Graph{
		Nodes:    make(map[int]struct{}),
		Edges:    make(map[int][]Edge),
		AllEdges: make([]Edge, 0),
		MaxNode:  -1,
	}

	// Start with two nodes connected by an edge (2 leaves)
	graph.AddNode(0)
	graph.AddNode(1)
	err := graph.AddEdge(0, 1, 1.0)
	if err != nil {
		return nil, err
	}

	// Keep adding leaves until we reach the desired number
	for countLeaves(graph) < numLeaves {
		err := addRandomLeaf(graph, rng)
		if err != nil {
			return nil, err
		}
	}

	return graph, nil
}

// Counts the number of leaf nodes (nodes with degree 1) in the graph
func countLeaves(graph *Graph) int {
	leafCount := 0
	for node := range graph.Nodes {
		if len(graph.Edges[node]) == 1 {
			leafCount++
		}
	}
	return leafCount
}

// Adds a new leaf to a random edge in the tree
func addRandomLeaf(graph *Graph, rng *rand.Rand) error {
	if len(graph.AllEdges) == 0 {
		return fmt.Errorf("cannot add leaf to empty graph")
	}

	// Select a random edge to split
	randomEdgeIndex := rng.Intn(len(graph.AllEdges))
	selectedEdge := graph.AllEdges[randomEdgeIndex]

	// Remove the selected edge
	_, err := graph.RemoveEdge(selectedEdge.Node1, selectedEdge.Node2)
	if err != nil {
		return err
	}

	// Create a new internal node
	newInternalNode := graph.AddNewNode()

	// Create a new leaf node
	newLeaf := graph.AddNewNode()

	// Connect the original nodes to the new internal node
	err = graph.AddEdge(selectedEdge.Node1, newInternalNode, 1.0)
	if err != nil {
		return err
	}

	err = graph.AddEdge(selectedEdge.Node2, newInternalNode, 1.0)
	if err != nil {
		return err
	}

	// Connect the new leaf to the new internal node
	err = graph.AddEdge(newInternalNode, newLeaf, 1.0)
	if err != nil {
		return err
	}

	return nil
}

// Returns a slice containing all leaf nodes in the graph
func GetLeafNodes(graph *Graph) []int {
	leaves := make([]int, 0)
	for node := range graph.Nodes {
		if len(graph.Edges[node]) == 1 {
			leaves = append(leaves, node)
		}
	}
	return leaves
}
