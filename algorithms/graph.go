package algorithms

import (
	"fmt"
	"math"
)

type Graph struct {
	Nodes map[int]struct{}
	Edges map[int][]Edge
	AllEdges []Edge
	MaxNode int
}

type Edge struct {
	Node1 int
	Node2 int
	Weight float64
}

func (e *Edge) SameAs(other *Edge) bool {
	return (e.Node1 == other.Node1 && e.Node2 == other.Node2) || (e.Node1 == other.Node2 && e.Node2 == other.Node1)
}

func IndexOfEdge[T Edge | *Edge](edges []T, start int, end int) int {
	for i, e := range edges {
		var edge *Edge
		switch v := any(e).(type) {
		case Edge:
			edge = &v
		case *Edge:
			edge = v
		}

		if edge.SameAs(&Edge{start, end, 0}) {
			return i
		}
	}

	return -1
}

func (g *Graph) AddNode(node int) bool {
	if _, ok := g.Nodes[node]; ok {
		return false
	}

	g.Nodes[node] = struct{}{}
	g.Edges[node] = make([]Edge, 0)

	if node > g.MaxNode {
		g.MaxNode = node
	}

	return true
}

func (g *Graph) AddNewNode() int {
	var node = g.MaxNode + 1
	g.AddNode(node)
	g.MaxNode = node
	return node
}

func (g *Graph) AddEdge(node1 int, node2 int, weight float64) error {
	if node1 > node2 {
		node1, node2 = node2, node1
	}

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
	if index := IndexOfEdge(g.AllEdges, node1, node2); index != -1 {
		return fmt.Errorf("edge %d-%d already exists", node1, node2)
	}

	var edge = Edge{node1, node2, weight}
	g.AllEdges = append(g.AllEdges, edge)
	g.Edges[node1] = append(g.Edges[node1], edge)
	g.Edges[node2] = append(g.Edges[node2], edge)

	return nil
}

func (g *Graph) RemoveEdge(node1 int, node2 int) (bool, error) {
	if node1 > node2 {
		node1, node2 = node2, node1
	}

	// fmt.Printf("Removing edge %d-%d\n", node1, node2)

	index1 := IndexOfEdge(g.Edges[node1], node1, node2)
	if index1 != -1 {
		g.Edges[node1] = append(g.Edges[node1][:index1], g.Edges[node1][index1+1:]...)
	}
	
	index2 := IndexOfEdge(g.Edges[node2], node1, node2)
	if index2 != -1 {
		g.Edges[node2] = append(g.Edges[node2][:index2], g.Edges[node2][index2+1:]...)
	}
	
	index3 := IndexOfEdge(g.AllEdges, node1, node2)
	if index3 != -1 {
		g.AllEdges = append(g.AllEdges[:index3], g.AllEdges[index3+1:]...)
	}

	if !(index1 != -1 && index2 != -1 && index3 != -1) && !(index1 == -1 && index2 == -1 && index3 == -1) {
		return false, fmt.Errorf("edge %d-%d found in only some lists: (from %d: %v, from %d: %v, from all: %v)", 
			node1, node2, node1, index1 != -1, node2, index2 != -1, index3 != -1)
	} 

	return index1 != -1 || index2 != -1 || index3 != -1, nil
}

func (g *Graph) MergeNodes(node1 int, node2 int) error {
	// fmt.Printf("Merging nodes %d and %d\n", node1, node2)

	if _, ok := g.Nodes[node1]; !ok {
		return fmt.Errorf("node %d does not exist in the graph", node1)
	}
	if _, ok := g.Nodes[node2]; !ok {
		return fmt.Errorf("node %d does not exist in the graph", node2)
	}

	if removed, err := g.RemoveEdge(node1, node2); err != nil {
		return err
	} else if !removed {
		return fmt.Errorf("edge %d-%d not found in the graph: merged nodes must be connected", node1, node2)
	}

	var edgesToChange []Edge
	for _, edge := range g.AllEdges {
		if edge.Node1 == node2 || edge.Node2 == node2 {
			edgesToChange = append(edgesToChange, edge)
		}
	}
	
	for _, edge := range edgesToChange {
		g.RemoveEdge(edge.Node1, edge.Node2)
		if edge.Node1 == node2 {
			g.AddEdge(node1, edge.Node2, edge.Weight)
		}
		if edge.Node2 == node2 {
			g.AddEdge(edge.Node1, node1, edge.Weight)
		}
	}
	
	delete(g.Nodes, node2)
	delete(g.Edges, node2)

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

		// PrintTree(g)
		err = g.ValidateTree()
		if err != nil {
			return err
		}

		merged[edge.Node2] = merged[edge.Node1]
	}

	return nil
}

func (g *Graph) SplitEdge(edge Edge, epsilon float64) error {
	removed, err := g.RemoveEdge(edge.Node1, edge.Node2)
	if err != nil {
		return err
	}
	if !removed {
		return fmt.Errorf("edge %d-%d not found in the graph", edge.Node1, edge.Node2)
	}
	
	var roundedWeight = (int)(math.Round(edge.Weight))
	if math.Abs(edge.Weight - float64(roundedWeight)) > epsilon {
		return fmt.Errorf("edge %d-%d has non-integer weight (%f)", edge.Node1, edge.Node2, edge.Weight)
	}

	var previousNode = edge.Node1
	for i := 0; i < roundedWeight-1; i++ {
		var newNode = g.AddNewNode()
		g.AddEdge(previousNode, newNode, 1)
		previousNode = newNode
	}

	g.AddEdge(previousNode, edge.Node2, 1)

	return nil
}

func (g *Graph) SplitEdges(epsilon float64) error {
	var edgesToSplit []Edge
	for _, edge := range g.AllEdges {
		if math.Abs(edge.Weight - 1) > epsilon {
			edgesToSplit = append(edgesToSplit, edge)
		}
	}

	for _, edge := range edgesToSplit {
		err := g.SplitEdge(edge, epsilon)
		if err != nil {
			return err
		}
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

func (g *Graph) ValidateTree() error {
	for _, edge := range g.AllEdges {
		if edge.Weight < 0 {
			return fmt.Errorf("edge %d-%d has negative weight", edge.Node1, edge.Node2)
		}
	}

	for _, edge := range g.AllEdges {
		if index1 := IndexOfEdge(g.Edges[edge.Node1], edge.Node1, edge.Node2); index1 == -1 {
			return fmt.Errorf("edge %d-%d not found in g.Edges[%d]: %v", edge.Node1, edge.Node2, edge.Node1, g.AllEdges)
		}
		if index2 := IndexOfEdge(g.Edges[edge.Node2], edge.Node1, edge.Node2); index2 == -1 {
			return fmt.Errorf("edge %d-%d not found in g.Edges[%d]: %v", edge.Node1, edge.Node2, edge.Node2, g.AllEdges)
		}
	}

	return nil
}
