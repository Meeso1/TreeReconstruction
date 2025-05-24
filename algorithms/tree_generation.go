package algorithms

import (
	"fmt"
	"math/rand"
)

// Creates a random tree with the specified number of leaves.
// The function starts with a simple path and randomly adds new leaves until
// the desired number is reached.
// chainExtensionProb is the probability of extending a chain by one additional node.
// connectToExistingProb is the probability of connecting new leaves to existing non-leaf nodes
// instead of splitting edges.
func GenerateRandomTree(
	numLeaves int,
	seed int64,
	chainExtensionProb float64,
	connectToExistingProb float64,
) (*Graph, error) {
	if numLeaves < 2 {
		return nil, fmt.Errorf("number of leaves must be at least 2, got %d", numLeaves)
	}

	if chainExtensionProb < 0 || chainExtensionProb >= 1 {
		return nil, fmt.Errorf("chainExtensionProb must be in range [0, 1), got %f", chainExtensionProb)
	}

	if connectToExistingProb < 0 || connectToExistingProb >= 1 {
		return nil, fmt.Errorf("connectToExistingProb must be in range [0, 1), got %f", connectToExistingProb)
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
		err := addRandomLeaf(graph, rng, chainExtensionProb, connectToExistingProb)
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
func addRandomLeaf(graph *Graph, rng *rand.Rand, chainExtensionProb float64, connectToExistingProb float64) error {
	if len(graph.AllEdges) == 0 {
		return fmt.Errorf("cannot add leaf to empty graph")
	}

	// Decide whether to connect to existing node or split an edge
	if rng.Float64() < connectToExistingProb {
		return addRandomLeafByConnecting(graph, rng, chainExtensionProb)
	} else {
		return addRandomLeafBySplitting(graph, rng, chainExtensionProb)
	}
}

// Adds a new leaf by connecting it to an existing non-leaf node
func addRandomLeafByConnecting(graph *Graph, rng *rand.Rand, chainExtensionProb float64) error {
	// Find all non-leaf nodes (nodes with degree > 1)
	nonLeafNodes := make([]int, 0)
	for node := range graph.Nodes {
		if len(graph.Edges[node]) > 1 {
			nonLeafNodes = append(nonLeafNodes, node)
		}
	}

	if len(nonLeafNodes) == 0 {
		// If no non-leaf nodes exist, fall back to splitting an edge
		return addRandomLeafBySplitting(graph, rng, chainExtensionProb)
	}

	// Select a random non-leaf node
	selectedNode := nonLeafNodes[rng.Intn(len(nonLeafNodes))]

	// Determine chain length
	chainLength := determineChainLength(rng, chainExtensionProb)

	// Create chain of nodes
	chainNodes := make([]int, chainLength)
	for i := 0; i < chainLength; i++ {
		chainNodes[i] = graph.AddNewNode()
	}

	// Connect the chain starting from the selected node
	err := graph.AddEdge(selectedNode, chainNodes[0], 1.0)
	if err != nil {
		return err
	}

	// Connect the chain nodes
	for i := 0; i < chainLength-1; i++ {
		err = graph.AddEdge(chainNodes[i], chainNodes[i+1], 1.0)
		if err != nil {
			return err
		}
	}

	// The last node in the chain is the leaf (chainNodes[chainLength-1])

	return nil
}

// Adds a new leaf by splitting a random edge in the tree
func addRandomLeafBySplitting(graph *Graph, rng *rand.Rand, chainExtensionProb float64) error {
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

	// Determine chain length
	chainLength := determineChainLength(rng, chainExtensionProb)

	// Create chain of nodes
	chainNodes := make([]int, chainLength)
	for i := 0; i < chainLength; i++ {
		chainNodes[i] = graph.AddNewNode()
	}

	// Connect the original nodes to the new internal node
	err = graph.AddEdge(selectedEdge.Node1, newInternalNode, 1.0)
	if err != nil {
		return err
	}

	err = graph.AddEdge(selectedEdge.Node2, newInternalNode, 1.0)
	if err != nil {
		return err
	}

	// Connect the chain starting from the new internal node
	err = graph.AddEdge(newInternalNode, chainNodes[0], 1.0)
	if err != nil {
		return err
	}

	// Connect the chain nodes
	for i := 0; i < chainLength-1; i++ {
		err = graph.AddEdge(chainNodes[i], chainNodes[i+1], 1.0)
		if err != nil {
			return err
		}
	}

	return nil
}

// Determines the chain length based on the extension probability
// Returns at least 1, and extends with probability chainExtensionProb up to a maximum length
func determineChainLength(rng *rand.Rand, chainExtensionProb float64) int {
	if chainExtensionProb <= 0 {
		return 1
	}

	// Calculate maximum chain length as 50 * chainExtensionProb, but at least 1
	maxLength := max(int(50*chainExtensionProb), 1)

	length := 1
	for length < maxLength && rng.Float64() < chainExtensionProb {
		length++
	}

	return length
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
