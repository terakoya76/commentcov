package pluggable

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/terakoya76/commentcov/pkg/common"
)

// PluginConfig represents commentcov plugin info.
type PluginConfig struct {
	Extension      string
	InstallCommand string `mapstructure:"install_command"`
	ExecuteCommand string `mapstructure:"execute_command"`
}

// FindPluginConfig returns the first matched plugin's PluginConfig by the extension.
func FindPluginConfig(pcs []PluginConfig, ext string) *PluginConfig {
	for _, pc := range pcs {
		if pc.Extension == ext {
			return &pc
		}
	}

	return nil
}

// Name returns the name.
func (pc *PluginConfig) Name() string {
	return fmt.Sprintf(
		"%s-plugin-for-%s",
		common.ProjectName,
		strings.ReplaceAll(pc.Extension, ".", ""),
	)
}

// Install installs the plugin for setup.
func (pc *PluginConfig) Install() error {
	cmd := pc.GetInstallCommand()
	return cmd.Run()
}

// GetInstallCommand returns *exec.Cmd to install the plugin.
func (pc *PluginConfig) GetInstallCommand() *exec.Cmd {
	cmds := strings.Split(pc.InstallCommand, " ")
	if len(cmds) == 1 {
		return exec.Command(pc.InstallCommand)
	}

	return exec.Command(cmds[0], cmds[1:]...)
}

// GetExecuteCommand returns *exec.Cmd to execute the plugin.
func (pc *PluginConfig) GetExecuteCommand() *exec.Cmd {
	cmds := strings.Split(pc.ExecuteCommand, " ")
	if len(cmds) == 1 {
		return exec.Command(pc.ExecuteCommand)
	}

	return exec.Command(cmds[0], cmds[1:]...)
}
