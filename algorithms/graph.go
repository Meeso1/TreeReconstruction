package algorithms

import (
	"fmt"
	"math"
)

type Graph struct {
	Nodes map[int]struct{}
	Edges map[int][]*Edge
	AllEdges []Edge
}

type Edge struct {
	Node1 int
	Node2 int
	Weight float64
}

func (g *Graph) AddNode(node int) bool {
	if _, ok := g.Nodes[node]; ok {
		return false
	}

	g.Nodes[node] = struct{}{}
	g.Edges[node] = make([]*Edge, 0)
	return true
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
	if weight < 0 {
		return fmt.Errorf("weight must be non-negative (got %f for edge %d-%d)", weight, node1, node2)
	}
	for _, edge := range g.Edges[node1] {
		if edge.Node2 == node2 {
			return fmt.Errorf("edge %d-%d already exists", node1, node2)
		}
	}

	g.AllEdges = append(g.AllEdges, Edge{node1, node2, weight})
	g.Edges[node1] = append(g.Edges[node1], &g.AllEdges[len(g.AllEdges) - 1])
	g.Edges[node2] = append(g.Edges[node2], &g.AllEdges[len(g.AllEdges) - 1])

	return nil
}

func (g *Graph) RemoveEdge(node1 int, node2 int) {
	for i, edge := range g.Edges[node1] {
		if edge.Node2 == node2 {
			g.Edges[node1] = append(g.Edges[node1][:i], g.Edges[node1][i+1:]...)
			break
		}
	}

	for i, edge := range g.Edges[node2] {
		if edge.Node1 == node1 {
			g.Edges[node2] = append(g.Edges[node2][:i], g.Edges[node2][i+1:]...)
			break
		}
	}

	for i, edge := range g.AllEdges {
		if (edge.Node1 == node1 && edge.Node2 == node2) {
			g.AllEdges = append(g.AllEdges[:i], g.AllEdges[i+1:]...)
			break
		}
	}
}

func (g *Graph) MergeNodes(node1 int, node2 int) error {
	if _, ok := g.Nodes[node1]; !ok {
		return fmt.Errorf("node %d does not exist in the graph", node1)
	}
	if _, ok := g.Nodes[node2]; !ok {
		return fmt.Errorf("node %d does not exist in the graph", node2)
	}

	delete(g.Nodes, node2)

	for _, edge := range g.Edges[node2] {
		if edge.Node1 == node1 || edge.Node2 == node1 {
			continue
		}

		if edge.Node1 == node2 {
			edge.Node1 = node1
		}
		if edge.Node2 == node2 {
			edge.Node2 = node1
		}
	}

	g.RemoveEdge(node1, node2)

	return nil
}

func (g *Graph) MergeZeroEdges(epsilon float64) error {
	var zeroEdges []Edge
	for _, edge := range g.AllEdges {
		if edge.Weight <= epsilon {
			zeroEdges = append(zeroEdges, edge)
		}
	}

	var merged = make(map[int]int)
	for _, edge := range zeroEdges {
		merged[edge.Node1] = edge.Node1
		merged[edge.Node2] = edge.Node2
	}

	for _, edge := range zeroEdges {
		err := g.MergeNodes(merged[edge.Node1], merged[edge.Node2])
		if err != nil {
			return err
		}

		merged[edge.Node2] = merged[edge.Node1]
	}

	return nil
}

func (g *Graph) IsIntegerWeighted(epsilon float64) bool {
	for _, edge := range g.AllEdges {
		if math.Abs(edge.Weight - math.Round(edge.Weight)) > epsilon {
			return false
		}
	}

	return true
}