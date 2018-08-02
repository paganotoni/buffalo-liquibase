package cmd

import (
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Allows to destroy generated code.",
	Aliases: []string{"d"},
}

// liquibaseCmd represents the liquibase command
var liquibaseCmd = &cobra.Command{
	Use:     "liquibase",
	Short:   "description about this plugin",
	Aliases: []string{"lb"},
}

func init() {
	rootCmd.AddCommand(liquibaseCmd)
}
