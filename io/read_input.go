package io

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseMatrix(fileContent string) ([][]uint32, error) {
	lines := strings.Split(fileContent, "\n")
	matrix := [][]uint32{}

	for _, line := range lines {
		fields := strings.Split(line, ",")
		row := make([]uint32, len(fields))
		for i, field := range fields {
			field = strings.TrimSpace(field)
			num, err := strconv.ParseUint(field, 10, 32)
			if err != nil {
				return nil, err
			}

			row[i] = uint32(num)
		}
		matrix = append(matrix, row)
	}


	if len(matrix) == 0 {
		return nil, errors.New("empty matrix")
	}

	if len(matrix[0]) == 0 {
		return nil, errors.New("empty matrix")
	}

	if len(matrix) != len(matrix[0]) {
		return nil, errors.New("matrix is not square")
	}

	for i, row := range matrix {
		if len(row) != len(matrix) {
			return nil, fmt.Errorf("row %d has %d elements, but row 0 has %d", i, len(row), len(matrix[0]))
		}
	}

	return matrix, nil
}
