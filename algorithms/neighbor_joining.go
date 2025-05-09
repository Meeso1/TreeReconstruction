package algorithms

import (
	"fmt"
	"math"
)

func MatrixToDict(matrix [][]float64) map[int]map[int]float64 {
	var dict = map[int]map[int]float64{}
	for i := range matrix {
		dict[i] = map[int]float64{}
		for j := range matrix[i] {
			dict[i][j] = matrix[i][j]
		}
	}
	return dict
}

func MakeRValues(distances map[int]map[int]float64, joinable map[int]struct{}) map[int]float64 {
	var r = make(map[int]float64)
	for i := range joinable {
		var sum = 0.0
		for j := range joinable {
			sum += distances[i][j]
		}
		r[i] = sum
	}

	return r
}

func NeighborJoining(matrix [][]float64) (*Graph, error) {
	for len(matrix) < 2 {
		return nil, fmt.Errorf("matrix must have at least 2 rows")
	}

	var joinable = map[int]struct{}{}
	for i := range matrix {
		joinable[i] = struct{}{}
	}
	var firstFreeNodeIndex = len(joinable)

	var distances = MatrixToDict(matrix)
	var tree = Graph{
		Nodes: map[int]struct{}{},
		Edges: map[int][]*Edge{},
	}

	for len(joinable) > 2 {
		var r = MakeRValues(distances, joinable)
		
		var minScore = math.Inf(1)
		var minI, minJ = -1, -1
		for i := range joinable {
			for j := range joinable {
				if i >= j {
					continue
				}
				
				var q = (float64(len(joinable)) - 2) * distances[i][j] - r[i] - r[j]
				if q < minScore {
					minScore = q
					minI = i
					minJ = j
				}
			}
		}			
		
		var u = firstFreeNodeIndex
		firstFreeNodeIndex++

		var distanceToI = (distances[minI][minJ] + (r[minI] - r[minJ]) / float64(len(joinable) - 2)) / 2
		var distanceToJ = distances[minI][minJ] - distanceToI

		//fmt.Printf("Adding nodes: u=%d, i=%d, j=%d\n", u, minI, minJ)
		//fmt.Printf("Distance to I: %f, Distance to J: %f\n", distanceToI, distanceToJ)

		tree.AddNode(u)
		tree.AddNode(minI)
		tree.AddNode(minJ)
		err := tree.AddEdge(minI, u, distanceToI)
		if err != nil {
			return nil, err
		}
		err = tree.AddEdge(minJ, u, distanceToJ)
		if err != nil {
			return nil, err
		}

		delete(joinable, minI)
		delete(joinable, minJ)
		joinable[u] = struct{}{}
		
		distances[u] = map[int]float64{}
		for k := range joinable {
			if k == u {
				distances[u][u] = 0
				continue
			}

			distances[u][k] = (distances[minI][k] + distances[minJ][k] - distances[minI][minJ]) / 2
			distances[k][u] = distances[u][k]
		}
		
		delete(distances, minI)
		delete(distances, minJ)

		for _, v := range distances {
			delete(v, minI)
			delete(v, minJ)
		}

		//fmt.Printf("Joinable: %v\n", joinable)
		//fmt.Printf("Distances: %v\n", distances)
	}

	var remaining []int
	for k := range joinable {
		remaining = append(remaining, k)
	}

	tree.AddNode(remaining[0])
	tree.AddNode(remaining[1])
	err := tree.AddEdge(remaining[0], remaining[1], distances[remaining[0]][remaining[1]])
	if err != nil {
		return nil, err
	}

	return &tree, nil
}
