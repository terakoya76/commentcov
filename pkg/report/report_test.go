package report_test

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov/pkg/report"
	"github.com/commentcov/commentcov/proto"
)

var (
	sortedSliceTrans = cmp.Transformer("Sort", func(in []string) []string {
		out := append([]string(nil), in...) // Copy input to avoid mutating it
		sort.Strings(out)
		return out
	})
)

func TestStringToMode(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want report.Mode
	}{

		{
			name: "file",
			str:  "file",
			want: report.ModeFile,
		},

		{
			name: "scope",
			str:  "scope",
			want: report.ModeScope,
		},

		{
			name: "file_scope",
			str:  "file_scope",
			want: report.ModeFileScope,
		},

		{
			name: "empty",
			str:  "",
			want: report.ModeInvalid,
		},

		{
			name: "random",
			str:  "hogehoge",
			want: report.ModeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := report.StringToMode(tt.str)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("*Counter values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestProfile(t *testing.T) {
	tests := []struct {
		name        string
		items       []*proto.CoverageItem
		wantCounter map[string]report.ScopedCounter
		wantFiles   []string
		wantScopes  []string
	}{
		{
			name: "single item",
			items: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   1,
						StartColumn: 1,
						EndLine:     10,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyFunction",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   1,
								StartColumn: 1,
								EndLine:     1,
								EndColumn:   18,
							},
							Comment: "MyFunction Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
			wantCounter: map[string]report.ScopedCounter{
				"hoge.go": {
					"PUBLIC_FUNCTION": &report.Counter{
						Covered: 1,
						Total:   1,
					},
				},
			},
			wantFiles: []string{
				"hoge.go",
			},
			wantScopes: []string{
				"PUBLIC_FUNCTION",
			},
		},

		{
			name: "multi items",
			items: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   1,
						StartColumn: 1,
						EndLine:     10,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyFunction",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   1,
								StartColumn: 1,
								EndLine:     1,
								EndColumn:   18,
							},
							Comment: "MyFunction Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   11,
						StartColumn: 11,
						EndLine:     20,
						EndColumn:   20,
					},
					File:       "hoge.go",
					Identifier: "MyFunction",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 11,
								EndLine:     11,
								EndColumn:   18,
							},
							Comment: "MyFunction Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
			wantCounter: map[string]report.ScopedCounter{
				"hoge.go": {
					"PUBLIC_FUNCTION": &report.Counter{
						Covered: 2,
						Total:   2,
					},
				},
			},
			wantFiles: []string{
				"hoge.go",
			},
			wantScopes: []string{
				"PUBLIC_FUNCTION",
			},
		},

		{
			name: "file duplicated",
			items: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   1,
						StartColumn: 1,
						EndLine:     10,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyFunction",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   1,
								StartColumn: 1,
								EndLine:     1,
								EndColumn:   18,
							},
							Comment: "MyFunction Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   11,
						StartColumn: 11,
						EndLine:     20,
						EndColumn:   20,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 11,
								EndLine:     11,
								EndColumn:   18,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
			wantCounter: map[string]report.ScopedCounter{
				"hoge.go": {
					"PUBLIC_FUNCTION": &report.Counter{
						Covered: 1,
						Total:   1,
					},
					"PUBLIC_VARIABLE": &report.Counter{
						Covered: 1,
						Total:   1,
					},
				},
			},
			wantFiles: []string{
				"hoge.go",
			},
			wantScopes: []string{
				"PUBLIC_FUNCTION",
				"PUBLIC_VARIABLE",
			},
		},

		{
			name: "scope duplicated",
			items: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   1,
						StartColumn: 1,
						EndLine:     10,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyFunction",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   1,
								StartColumn: 1,
								EndLine:     1,
								EndColumn:   18,
							},
							Comment: "MyFunction Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   11,
						StartColumn: 11,
						EndLine:     20,
						EndColumn:   20,
					},
					File:       "fuga.go",
					Identifier: "MyFunction",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 11,
								EndLine:     11,
								EndColumn:   18,
							},
							Comment: "MyFunction Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
			wantCounter: map[string]report.ScopedCounter{
				"hoge.go": {
					"PUBLIC_FUNCTION": &report.Counter{
						Covered: 1,
						Total:   1,
					},
				},
				"fuga.go": {
					"PUBLIC_FUNCTION": &report.Counter{
						Covered: 1,
						Total:   1,
					},
				},
			},
			wantFiles: []string{
				"hoge.go",
				"fuga.go",
			},
			wantScopes: []string{
				"PUBLIC_FUNCTION",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCounter, gotFiles, gotScopes := report.Profile(tt.items)

			if diff := cmp.Diff(tt.wantCounter, gotCounter); diff != "" {
				t.Errorf("map[string]report.ScopedCounter values are mismatch (-want +got):%s\n", diff)
			}

			if diff := cmp.Diff(tt.wantFiles, gotFiles, sortedSliceTrans); diff != "" {
				t.Errorf("[]string values are mismatch (-want +got):%s\n", diff)
			}

			if diff := cmp.Diff(tt.wantScopes, gotScopes, sortedSliceTrans); diff != "" {
				t.Errorf("[]string values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}
