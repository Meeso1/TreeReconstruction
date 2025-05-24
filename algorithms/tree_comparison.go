package algorithms

import (
	"sort"
)

// Checks if two trees have the same topology (structure)
// ignoring node numbering/names
func CompareTreeTopology(tree1, tree2 *Graph) bool {
	if len(tree1.Nodes) != len(tree2.Nodes) {
		return false
	}

	if len(tree1.AllEdges) != len(tree2.AllEdges) {
		return false
	}

	if !compareDegreeSequences(tree1, tree2) {
		return false
	}

	canonical1 := generateCanonicalRepresentation(tree1)
	canonical2 := generateCanonicalRepresentation(tree2)

	return canonical1 == canonical2
}

// Checks if two trees have the same degree sequence
func compareDegreeSequences(tree1, tree2 *Graph) bool {
	degrees1 := make([]int, 0)
	degrees2 := make([]int, 0)

	for node := range tree1.Nodes {
		degrees1 = append(degrees1, len(tree1.Edges[node]))
	}

	for node := range tree2.Nodes {
		degrees2 = append(degrees2, len(tree2.Edges[node]))
	}

	sort.Ints(degrees1)
	sort.Ints(degrees2)

	if len(degrees1) != len(degrees2) {
		return false
	}

	for i := range degrees1 {
		if degrees1[i] != degrees2[i] {
			return false
		}
	}

	return true
}

// Creates a canonical string representation of the tree
func generateCanonicalRepresentation(tree *Graph) string {
	if len(tree.Nodes) == 0 {
		return ""
	}

	// Find tree centers (there can be 1 or 2 centers in a tree)
	centers := findTreeCenters(tree)

	// Generate canonical representation from each center and take the lexicographically smallest
	var minRepresentation string

	for _, center := range centers {
		representation := generateRepresentationFromRoot(tree, center)
		if minRepresentation == "" || representation < minRepresentation {
			minRepresentation = representation
		}
	}

	return minRepresentation
}

// Finds the center(s) of the tree using the standard algorithm
func findTreeCenters(tree *Graph) []int {
	if len(tree.Nodes) == 1 {
		for node := range tree.Nodes {
			return []int{node}
		}
	}

	// Start with all nodes
	remaining := make(map[int]bool)
	degrees := make(map[int]int)

	for node := range tree.Nodes {
		remaining[node] = true
		degrees[node] = len(tree.Edges[node])
	}

	// Repeatedly remove leaves until we have 1 or 2 nodes left
	for len(remaining) > 2 {
		// Find all current leaves (degree 1)
		leaves := make([]int, 0)
		for node := range remaining {
			if degrees[node] == 1 {
				leaves = append(leaves, node)
			}
		}

		// Remove leaves and update degrees of their neighbors
		for _, leaf := range leaves {
			delete(remaining, leaf)

			// Find the neighbor of this leaf and decrease its degree
			for _, edge := range tree.Edges[leaf] {
				var neighbor int
				if edge.Node1 == leaf {
					neighbor = edge.Node2
				} else {
					neighbor = edge.Node1
				}

				if remaining[neighbor] {
					degrees[neighbor]--
				}
			}
		}
	}

	// Return remaining nodes as centers
	centers := make([]int, 0)
	for node := range remaining {
		centers = append(centers, node)
	}
	sort.Ints(centers) // For deterministic order
	return centers
}

// Creates a canonical string representation rooted at the given node
func generateRepresentationFromRoot(tree *Graph, root int) string {
	visited := make(map[int]bool)
	return dfsCanonical(tree, root, visited)
}

// Performs DFS to generate canonical representation
func dfsCanonical(tree *Graph, node int, visited map[int]bool) string {
	visited[node] = true

	// Get all unvisited neighbors
	childRepresentations := make([]string, 0)

	for _, edge := range tree.Edges[node] {
		var neighbor int
		if edge.Node1 == node {
			neighbor = edge.Node2
		} else {
			neighbor = edge.Node1
		}

		if !visited[neighbor] {
			childRepresentation := dfsCanonical(tree, neighbor, visited)
			childRepresentations = append(childRepresentations, childRepresentation)
		}
	}

	// Sort child representations for canonical order
	sort.Strings(childRepresentations)

	// Create representation: (child1)(child2)...(childN)
	result := "("
	for _, childRepr := range childRepresentations {
		result += childRepr
	}
	result += ")"

	return result
}
