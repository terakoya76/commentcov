package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/commentcov/commentcov/pkg/execute"
)

// coverageCmd returns coverage info.
var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "generate coverage reports",
	Long:  "generate coverage reports",
	Run: func(_ *cobra.Command, _ []string) {
		if err := execute.Run(cfg); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	},
}
