package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/terakoya76/commentcov/pkg/common"
)

// versionCmd returns the version info.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s", common.Version),
	Long:  fmt.Sprintf(`All software has versions. This is %s's`, common.Version),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", common.Version)
		os.Exit(0)
	},
}
