package config

import (
	"github.com/terakoya76/commentcov/pkg/pluggable"
)

// ViperConfig is unmarshaled from viper.Config.
type ViperConfig struct {
	TargetPath   string   `mapstructure:"target_path"`
	ExcludePaths []string `mapstructure:"exclude_paths"`
	Plugins      []pluggable.PluginConfig
	Mode         string
}
