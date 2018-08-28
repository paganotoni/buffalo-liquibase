package cmd

import (
	"github.com/spf13/cobra"
)

// liquibaseCmd represents the liquibase command
var liquibaseCmd = &cobra.Command{
	Use:     "liquibase",
	Short:   "Runs yours migrations using liquibase",
	Aliases: []string{"lb"},
}

func init() {
	rootCmd.AddCommand(liquibaseCmd)
}
