package cmd

import (
	"fmt"

	"github.com/paganotoni/buffalo-liquibase/liquibase"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of liquibase",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("liquibase", liquibase.Version)
		return nil
	},
}

func init() {
	liquibaseCmd.AddCommand(versionCmd)
}
