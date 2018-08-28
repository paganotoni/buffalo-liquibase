package cmd

import "github.com/spf13/cobra"

// upCmd runs /migrations or --path up against buffalo db with liquibase
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return upCmd.RunE(cmd, args)
	},
}

func init() {
	liquibaseCmd.AddCommand(migrateCmd)
}
