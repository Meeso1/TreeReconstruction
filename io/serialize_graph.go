package io

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"treereconstruction/algorithms"
)

type SerializationType int

const (
	SerializationTypeBrackets SerializationType = iota
	SerializationTypeBracketsShortened
	SerializationTypeNeighborLists
)

func MakePrefixSuffix(incomingEdgeLength int, useShortenedSyntax bool) (string, string) {
	if !useShortenedSyntax || incomingEdgeLength == 1 {
		return strings.Repeat("(", incomingEdgeLength), strings.Repeat(")", incomingEdgeLength)
	}
	return "[" + strconv.Itoa(incomingEdgeLength) + "](", ")"
}

func SerializeChildrenAsBrackets(
	graph *algorithms.Graph,
	node int,
	alreadySerialized *map[int]struct{},
	incomingEdgeLength int,
	useShortenedSyntax bool,
) (string, error) {
	if _, ok := (*alreadySerialized)[node]; ok {
		return "", nil
	}

	(*alreadySerialized)[node] = struct{}{}
	var result, suffix = MakePrefixSuffix(incomingEdgeLength, useShortenedSyntax)
	for _, edge := range graph.Edges[node] {
		var otherNode = edge.Node1
		if otherNode == node {
			otherNode = edge.Node2
		}

		if _, ok := (*alreadySerialized)[otherNode]; ok {
			continue
		}

		var inside, err = SerializeChildrenAsBrackets(graph, otherNode, alreadySerialized, int(math.Round(edge.Weight)), useShortenedSyntax)
		if err != nil {
			return "", err
		}

		result += inside
	}
	return result + suffix, nil
}

func SerializeChildrenAsNeighborLists(graph *algorithms.Graph) (string, error) {
	var allNodes = make([]int, 0)
	for node := range graph.Nodes {
		allNodes = append(allNodes, node)
	}
	sort.Ints(allNodes)
	
	var result = ""
	for _, node := range allNodes {
		result += fmt.Sprintf("%d:", node)

		var edges = graph.Edges[node]

		var neighbors = make([]int, 0)
		for _, edge := range edges {
			var otherNode = edge.Node1
			if otherNode == node {
				otherNode = edge.Node2
			}
			neighbors = append(neighbors, otherNode)
		}
		sort.Ints(neighbors)

		if len(neighbors) == 0 {
			return "", fmt.Errorf("node %d has no neighbors", node)
		}

		for _, neighbor := range neighbors[:len(neighbors)-1] {
			result += fmt.Sprintf("%d,", neighbor)
		}
		result += fmt.Sprintf("%d;\n", neighbors[len(neighbors)-1])
	}

	return result, nil
}

func SerializeGraph(graph *algorithms.Graph, serializationType SerializationType) (string, error) {
	switch serializationType {
	case SerializationTypeBrackets:
		return SerializeChildrenAsBrackets(graph, 0, &map[int]struct{}{}, 1, false)
	case SerializationTypeBracketsShortened:
		return SerializeChildrenAsBrackets(graph, 0, &map[int]struct{}{}, 1, true)
	case SerializationTypeNeighborLists:
		err := graph.SplitEdges(1e-6)
		if err != nil {
			return "", err
		}
		return SerializeChildrenAsNeighborLists(graph)
	default:
		return "", fmt.Errorf("invalid serialization type: %d", serializationType)
	}
}
