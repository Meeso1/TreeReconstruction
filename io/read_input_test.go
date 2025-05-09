package io

import (
	"testing"
)

func MatrixEquals(a, b [][]uint32) bool {
	if len(a) != len(b) {
		return false
	}
	
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func TestParseMatrix(t *testing.T) {
	tests := []struct {
		name string
		input string
		want [][]uint32
		wantErr bool
	}{
		{
			name: "valid matrix",
			input: "1,2,3\n2,2,3\n3,2,1",
			want: [][]uint32{{1, 2, 3}, {2, 2, 3}, {3, 2, 1}},
			wantErr: false,
		},
		{
			name: "matrix with spaces",
			input: "1, 2, 3\n2, 2, 3\n3, 2, 1",
			want: [][]uint32{{1, 2, 3}, {2, 2, 3}, {3, 2, 1}},
			wantErr: false,
		},
		{
			name: "empty matrix",
			input: "",
			want: nil,
			wantErr: true,
		},
		{
			name: "matrix with empty row",
			input: "\n",
			want: nil,
			wantErr: true,
		},
		{
			name: "non-square matrix",
			input: "1,2,3\n2,2,3\n3,2",
			want: nil,
			wantErr: true,
		},
		{
			name: "matrix with non-integer values",
			input: "1,2,3\n2,2,3\n3,2,foo",
			want: nil,
			wantErr: true,
		},
		{
			name: "matrix with negative values",
			input: "1,2,3\n2,2,3\n3,2,-1",
			want: nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseMatrix(test.input)
			if err != nil && !test.wantErr {
				t.Errorf("ParseMatrix(%q) returned error: %v", test.input, err)
			} else if err == nil && test.wantErr {
				t.Errorf("ParseMatrix(%q) was expected to return an error, but got nil", test.input)
			}

			if !MatrixEquals(got, test.want) {
				t.Errorf("ParseMatrix(%q) returned incorrect matrix: expected %v, got %v", test.input, test.want, got)
			}
		})
	}
}
