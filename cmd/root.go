package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commentcov/commentcov/pkg/common"
	"github.com/commentcov/commentcov/pkg/config"
)

var (
	// cfgPath is the flag variable representing config filepath.
	cfgPath string
	// cfg holds info from config.
	cfg config.ViperConfig
)

// rootCmd returns help cmd.
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
	},
}

// Execute is the callable interface of cmd package.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}

// nolint: gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	rootCmd.PersistentFlags().StringVarP(
		&cfgPath,
		"config",
		"c",
		"",
		fmt.Sprintf("config file (default is $HOME/.%s.yaml)", common.ProjectName),
	)

	rootCmd.AddCommand(
		coverageCmd,
		versionCmd,
	)
}

func initConfig() {
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		cur, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		viper.AddConfigPath(fmt.Sprintf("%s/", cur))
		viper.SetConfigName(fmt.Sprintf(".%s", common.ProjectName))
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	// set default value
	cfg = config.ViperConfig{
		TargetPath: ".",
		Mode:       "file",
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}
