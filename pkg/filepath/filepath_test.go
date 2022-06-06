package filepath_test

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov/pkg/filepath"
)

var (
	sortedSliceTrans = cmp.Transformer("Sort", func(in []string) []string {
		out := append([]string(nil), in...) // Copy input to avoid mutating it
		sort.Strings(out)
		return out
	})
)

func TestExtensions(t *testing.T) {
	tests := []struct {
		name    string
		fileset *filepath.Fileset
		want    []string
	}{

		{
			name: "standard",
			fileset: &filepath.Fileset{
				Set: map[string]map[string]struct{}{
					".go": {
						"hoge.go": {},
						"fuga.go": {},
					},
					".rs": {
						"hoge.rs": {},
						"fuga.rs": {},
					},
				},
			},
			want: []string{
				".go",
				".rs",
			},
		},

		{
			name: "empty",
			fileset: &filepath.Fileset{
				Set: map[string]map[string]struct{}{},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileset.Extensions()
			if diff := cmp.Diff(tt.want, got, sortedSliceTrans); diff != "" {
				t.Errorf("[]string values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestFiles(t *testing.T) {
	tests := []struct {
		name      string
		fileset   *filepath.Fileset
		extension string
		want      []string
	}{

		{
			name: "standard",
			fileset: &filepath.Fileset{
				Set: map[string]map[string]struct{}{
					".go": {
						"hoge.go": {},
						"fuga.go": {},
					},
					".rs": {
						"hoge.rs": {},
						"fuga.rs": {},
					},
				},
			},
			extension: ".go",
			want: []string{
				"hoge.go",
				"fuga.go",
			},
		},

		{
			name: "empty",
			fileset: &filepath.Fileset{
				Set: map[string]map[string]struct{}{},
			},
			extension: ".go",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fileset.Files(tt.extension)
			if diff := cmp.Diff(tt.want, got, sortedSliceTrans); diff != "" {
				t.Errorf("[]string values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}
