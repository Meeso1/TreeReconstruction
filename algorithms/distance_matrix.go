package algorithms

import (
	"fmt"
	"sort"
)

// Computes the distance matrix between all leaf nodes in the tree.
// The distances are calculated using BFS to find shortest paths between leaves.
func CalculateDistanceMatrix(graph *Graph) ([][]int, error) {
	leaves := GetLeafNodes(graph)
	if len(leaves) == 0 {
		return nil, fmt.Errorf("no leaf nodes found in the graph")
	}

	// Sort leaves for consistent ordering
	sort.Ints(leaves)

	n := len(leaves)
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	// Create a mapping from node ID to matrix index
	leafToIndex := make(map[int]int)
	for i, leaf := range leaves {
		leafToIndex[leaf] = i
	}

	// Calculate distances between all pairs of leaves
	for i, leaf1 := range leaves {
		distances := bfsDistances(graph, leaf1)
		for j, leaf2 := range leaves {
			if distance, exists := distances[leaf2]; exists {
				matrix[i][j] = distance
			} else {
				return nil, fmt.Errorf("no path found between leaves %d and %d", leaf1, leaf2)
			}
		}
	}

	return matrix, nil
}

// Performs BFS from the start node and returns distances to all reachable nodes
func bfsDistances(graph *Graph, start int) map[int]int {
	distances := make(map[int]int)
	visited := make(map[int]bool)
	queue := []int{start}

	distances[start] = 0
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range graph.Edges[current] {
			var neighbor int
			if edge.Node1 == current {
				neighbor = edge.Node2
			} else {
				neighbor = edge.Node1
			}

			if !visited[neighbor] {
				visited[neighbor] = true
				distances[neighbor] = distances[current] + int(edge.Weight)
				queue = append(queue, neighbor)
			}
		}
	}

	return distances
}

// Converts a distance matrix to the CSV string format used by the application
func FormatDistanceMatrix(matrix [][]int) string {
	if len(matrix) == 0 {
		return ""
	}

	result := ""
	for i, row := range matrix {
		for j, val := range row {
			if j > 0 {
				result += ","
			}
			result += fmt.Sprintf("%d", val)
		}
		if i < len(matrix)-1 {
			result += "\n"
		}
	}

	return result
}
