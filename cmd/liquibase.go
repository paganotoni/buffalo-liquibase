package cmd

import (
	"github.com/spf13/cobra"
)

// liquibaseCmd represents the liquibase command
var liquibaseCmd = &cobra.Command{
	Use:     "liquibase",
	Short:   "description about this plugin",
	Aliases: []string{"lb"},
}

func init() {
	rootCmd.AddCommand(liquibaseCmd)
}
