package io

import (
	"math"
	"strconv"
	"strings"
	"treereconstruction/algorithms"
)

func MakePrefixSuffix(incomingEdgeLength int, useShortenedSyntax bool) (string, string) {
	if !useShortenedSyntax || incomingEdgeLength == 1 {
		return strings.Repeat("(", incomingEdgeLength), strings.Repeat(")", incomingEdgeLength)
	}
	return "[" + strconv.Itoa(incomingEdgeLength) + "](", ")"
}

func SerializeChildren(
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

		var inside, err = SerializeChildren(graph, otherNode, alreadySerialized, int(math.Round(edge.Weight)), useShortenedSyntax)
		if err != nil {
			return "", err
		}

		result += inside
	}
	return result + suffix, nil
}

func SerializeGraph(graph *algorithms.Graph, useShortenedSyntax bool) (string, error) {
	var result, err = SerializeChildren(graph, 0, &map[int]struct{}{}, 1, useShortenedSyntax)
	if err != nil {
		return "", err
	}
	return result, nil
}
