package algorithms

import "fmt"

type Graph struct {
	Nodes map[int]struct{}
	Edges [][]Edge
}

type Edge struct {
	Node1 int
	Node2 int
	Weight float64
}

func (g *Graph) AddNode(node int) {
	g.Nodes[node] = struct{}{}
	g.Edges = append(g.Edges, make([]Edge, 0))
}

func (g *Graph) AddEdge(node1 int, node2 int, weight float64) error {
	if _, ok := g.Nodes[node1]; !ok {
		return fmt.Errorf("node %d does not exist in the graph", node1)
	}
	if _, ok := g.Nodes[node2]; !ok {
		return fmt.Errorf("node %d does not exist in the graph", node2)
	}
	if node1 == node2 {
		return fmt.Errorf("node %d cannot be connected to itself", node1)
	}
	if weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}
	for _, edge := range g.Edges[node1] {
		if edge.Node2 == node2 {
			return fmt.Errorf("edge %d-%d already exists", node1, node2)
		}
	}

	var edge = Edge{node1, node2, weight}
	g.Edges[node1] = append(g.Edges[node1], edge)
	g.Edges[node2] = append(g.Edges[node2], edge)

	return nil
}
