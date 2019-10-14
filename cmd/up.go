package cmd

import (
	"github.com/gobuffalo/pop"
	// "github.com/paganotoni/buffalo-liquibase/liquibase/migrator"
	"github.com/spf13/cobra"
)

var changeLogFile string
var environment string
var databaseYmlFile string

// upCmd runs /migrations or --path up against buffalo db with liquibase
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "runs liquibase migrations up",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		//TODO: Do the actual work and run migrations UP.
		// migrator := migrator.PostgresMigrator{}
		// if err := migrator.Prepare(); err != nil {
		// 	return err
		// }

		return nil
	},
}

func init() {
	upCmd.PersistentFlags().StringVar(&changeLogFile, "c", "./migrations/changelog.xml", "migrations changelog")
	upCmd.PersistentFlags().StringVar(&environment, "e", "development", "environment to run the migrations against")
	migrateCmd.AddCommand(upCmd)
}
