package pluggable_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov/pkg/common"
	"github.com/commentcov/commentcov/pkg/pluggable"
)

func TestBatched(t *testing.T) {
	tests := []struct {
		name      string
		pcs       []pluggable.PluginConfig
		extension string
		want      *pluggable.PluginConfig
	}{

		{
			name: "standard",
			pcs: []pluggable.PluginConfig{
				{
					Extension:      ".go",
					InstallCommand: "hoge",
					ExecuteCommand: "fuga",
				},
				{
					Extension:      ".rs",
					InstallCommand: "hoge",
					ExecuteCommand: "fuga",
				},
			},
			extension: ".go",
			want: &pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "hoge",
				ExecuteCommand: "fuga",
			},
		},

		{
			name: "different extension",
			pcs: []pluggable.PluginConfig{
				{
					Extension:      ".go",
					InstallCommand: "hoge",
					ExecuteCommand: "fuga",
				},
				{
					Extension:      ".rs",
					InstallCommand: "hoge",
					ExecuteCommand: "fuga",
				},
			},
			extension: ".ts",
			want:      nil,
		},

		{
			name:      "empty",
			pcs:       []pluggable.PluginConfig{},
			extension: ".go",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pluggable.FindPluginConfig(tt.pcs, tt.extension)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("*PluginConfig values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		name string
		pc   pluggable.PluginConfig
		want string
	}{

		{
			name: "standard",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "hoge",
				ExecuteCommand: "fuga",
			},
			want: fmt.Sprintf("%s-plugin-for-go", common.ProjectName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pc.Name()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("string values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestGetInstallCommand(t *testing.T) {
	tests := []struct {
		name string
		pc   pluggable.PluginConfig
		want *exec.Cmd
	}{

		{
			name: "1 words command",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "ls",
				ExecuteCommand: "echo fuga",
			},
			want: exec.Command("ls"),
		},

		{
			name: "2 words command",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "echo hoge",
				ExecuteCommand: "echo fuga",
			},
			want: exec.Command("echo", "hoge"),
		},

		{
			name: "3 words command",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "git grep hoge",
				ExecuteCommand: "echo fuga",
			},
			want: exec.Command("git", "grep", "hoge"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pc.GetInstallCommand()
			if diff := cmp.Diff(tt.want.Path, got.Path); diff != "" {
				t.Errorf("*exec.Cmd values are mismatch (-want +got):%s\n", diff)
			}

			if diff := cmp.Diff(tt.want.Args, got.Args); diff != "" {
				t.Errorf("*exec.Cmd values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestGetExecuteCommand(t *testing.T) {
	tests := []struct {
		name string
		pc   pluggable.PluginConfig
		want *exec.Cmd
	}{

		{
			name: "1 words command",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "echo hoge",
				ExecuteCommand: "ls",
			},
			want: exec.Command("ls"),
		},

		{
			name: "2 words command",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "echo hoge",
				ExecuteCommand: "echo fuga",
			},
			want: exec.Command("echo", "fuga"),
		},

		{
			name: "3 words command",
			pc: pluggable.PluginConfig{
				Extension:      ".go",
				InstallCommand: "echo hoge",
				ExecuteCommand: "git grep fuga",
			},
			want: exec.Command("git", "grep", "fuga"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pc.GetExecuteCommand()
			if diff := cmp.Diff(tt.want.Path, got.Path); diff != "" {
				t.Errorf("*exec.Cmd values are mismatch (-want +got):%s\n", diff)
			}

			if diff := cmp.Diff(tt.want.Args, got.Args); diff != "" {
				t.Errorf("*exec.Cmd values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}
