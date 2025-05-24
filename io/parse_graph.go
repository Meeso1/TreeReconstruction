package io

import (
	"fmt"
	"strconv"
	"strings"
	"treereconstruction/algorithms"
)

// Parses a neighbor list format string into a Graph structure
func ParseNeighborList(content string) (*algorithms.Graph, error) {
	graph := &algorithms.Graph{
		Nodes:    make(map[int]struct{}),
		Edges:    make(map[int][]algorithms.Edge),
		AllEdges: make([]algorithms.Edge, 0),
		MaxNode:  -1,
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.Contains(line, ":") || !strings.HasSuffix(line, ";") {
			return nil, fmt.Errorf("invalid format: line must be 'node:neighbors;', got: %s", line)
		}

		line = strings.TrimSuffix(line, ";")
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format: expected 'node:neighbors', got: %s", line)
		}

		nodeStr := strings.TrimSpace(parts[0])
		node, err := strconv.Atoi(nodeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid node ID: %s", nodeStr)
		}

		graph.AddNode(node)
	}

	addedEdges := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = strings.TrimSuffix(line, ";")
		parts := strings.Split(line, ":")

		nodeStr := strings.TrimSpace(parts[0])
		node, _ := strconv.Atoi(nodeStr)

		neighborStr := strings.TrimSpace(parts[1])
		if neighborStr == "" {
			return nil, fmt.Errorf("node %d has no neighbors", node)
		}

		neighbors := strings.Split(neighborStr, ",")
		for _, neighborStr := range neighbors {
			neighborStr = strings.TrimSpace(neighborStr)
			neighbor, err := strconv.Atoi(neighborStr)
			if err != nil {
				return nil, fmt.Errorf("invalid neighbor ID: %s", neighborStr)
			}

			var edgeKey string
			if node < neighbor {
				edgeKey = fmt.Sprintf("%d-%d", node, neighbor)
			} else {
				edgeKey = fmt.Sprintf("%d-%d", neighbor, node)
			}

			if !addedEdges[edgeKey] {
				err := graph.AddEdge(node, neighbor, 1.0)
				if err != nil {
					return nil, fmt.Errorf("error adding edge %d-%d: %v", node, neighbor, err)
				}
				addedEdges[edgeKey] = true
			}
		}
	}

	return graph, nil
}
