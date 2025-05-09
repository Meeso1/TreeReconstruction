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
		for j := range distances {
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
		Edges: [][]*Edge{},
	}

	for len(joinable) > 2 {
		var r = MakeRValues(distances, joinable)
		
		var minScore = math.Inf(1)
		var minI, minJ = -1, -1
		for i := range joinable {
			for j := range joinable {
				if i == j {
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

		tree.AddNode(u)
		tree.AddEdge(minI, u, distanceToI)
		tree.AddEdge(minJ, u, distanceToJ)

		delete(joinable, minI)
		delete(joinable, minJ)
		joinable[u] = struct{}{}
		
		for k := range joinable {
			if k == u {
				continue
			}
			
			distances[u] = map[int]float64{}
			distances[u][k] = (distances[minI][k] + distances[minJ][k] - distances[minI][minJ]) / 2
			distances[k][u] = distances[u][k]
		}
		
		delete(distances, minI)
		delete(distances, minJ)
	}

	var remaining []int
	for k := range joinable {
		remaining = append(remaining, k)
	}

	tree.AddEdge(remaining[0], remaining[1], distances[remaining[0]][remaining[1]])

	return &tree, nil
}
