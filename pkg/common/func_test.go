package common_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov/pkg/common"
)

func TestBatched(t *testing.T) {
	tests := []struct {
		name      string
		slice     []string
		batchSize int
		want      [][]string
	}{

		{
			name: "standard",
			slice: []string{
				"a",
				"b",
				"c",
				"d",
			},
			batchSize: 2,
			want: [][]string{
				{"a", "b"},
				{"c", "d"},
			},
		},

		{
			name: "batchSize 1",
			slice: []string{
				"a",
				"b",
				"c",
				"d",
			},
			batchSize: 1,
			want: [][]string{
				{"a"},
				{"b"},
				{"c"},
				{"d"},
			},
		},

		{
			name: "same batchSize as slice length",
			slice: []string{
				"a",
				"b",
				"c",
				"d",
			},
			batchSize: 4,
			want: [][]string{
				{"a", "b", "c", "d"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := common.Batched(tt.slice, tt.batchSize)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("[][]string values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}
