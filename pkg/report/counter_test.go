package report_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov/pkg/report"
	"github.com/commentcov/commentcov/proto"
)

func TestNewCounter(t *testing.T) {
	tests := []struct {
		name string
		want *report.Counter
	}{

		{
			name: "standard",
			want: &report.Counter{
				Covered: 0,
				Total:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := report.NewCounter()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("*Counter values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestCalcRate(t *testing.T) {
	tests := []struct {
		name    string
		counter *report.Counter
		want    float64
	}{

		{
			name: "standard",
			counter: &report.Counter{
				Covered: 5,
				Total:   10,
			},
			want: 50.0,
		},

		{
			name: "100%",
			counter: &report.Counter{
				Covered: 10,
				Total:   10,
			},
			want: 100.0,
		},

		{
			name: "0%",
			counter: &report.Counter{
				Covered: 0,
				Total:   0,
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.counter.CalcRate()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("*Counter values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestProfile(t *testing.T) {
	tests := []struct {
		name    string
		counter *report.Counter
		item    *proto.CoverageItem
		want    *report.Counter
	}{

		{
			name: "standard",
			counter: &report.Counter{
				Covered: 0,
				Total:   0,
			},
			item: &proto.CoverageItem{
				HeaderComments: []*proto.Comment{
					{
						Block:   &proto.Block{},
						Comment: "Header1",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block:   &proto.Block{},
						Comment: "Inline1",
					},
				},
			},
			want: &report.Counter{
				Covered: 1,
				Total:   1,
			},
		},

		{
			name: "no header",
			counter: &report.Counter{
				Covered: 0,
				Total:   0,
			},
			item: &proto.CoverageItem{
				HeaderComments: []*proto.Comment{},
				InlineComments: []*proto.Comment{
					{
						Block:   &proto.Block{},
						Comment: "Inline1",
					},
				},
			},
			want: &report.Counter{
				Covered: 0,
				Total:   1,
			},
		},

		{
			name: "no inline",
			counter: &report.Counter{
				Covered: 0,
				Total:   0,
			},
			item: &proto.CoverageItem{
				HeaderComments: []*proto.Comment{
					{
						Block:   &proto.Block{},
						Comment: "Header1",
					},
				},
				InlineComments: []*proto.Comment{},
			},
			want: &report.Counter{
				Covered: 1,
				Total:   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.counter.Profile(tt.item)
			got := tt.counter
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("*Counter values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}
