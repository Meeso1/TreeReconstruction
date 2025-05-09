package algorithms

func CastMatrixToFloat(matrix [][]uint32) [][]float64 {
	var floatMatrix = make([][]float64, len(matrix))
	for i := range matrix {
		floatMatrix[i] = make([]float64, len(matrix[i]))
		for j := range matrix[i] {
			floatMatrix[i][j] = float64(matrix[i][j])
		}
	}

	return floatMatrix
}

func ReconstructIntTree(matrix [][]uint32, epsilon float64) (*Graph, error) {
	var floatMatrix = CastMatrixToFloat(matrix)
	var tree, err = NeighborJoining(floatMatrix)
	if err != nil {
		return nil, err
	}

	err = tree.MergeZeroEdges(epsilon)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

