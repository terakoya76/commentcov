package report_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov/pkg/report"
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
