package config

import (
	"github.com/terakoya76/commentcov/pkg/pluggable"
)

// ViperConfig is unmarshaled from viper.Config.
type ViperConfig struct {
	TargetPath    string   `mapstructure:"target_path"`
	ExcludePathes []string `mapstructure:"exclude_pathes"`
	Plugins       []pluggable.PluginConfig
	Mode          string
}
