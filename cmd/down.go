package cmd

import (
	"errors"
	"os/exec"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var rollbackCount int

// downCmd runs /migrations down
var downCmd = &cobra.Command{
	Use:     "rollback",
	Aliases: []string{"down"},
	Short:   "rollbacks passed number of migrations back, otherwise it runs one migration back",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("liquibase"); err != nil {
			return errors.New("could not find liquibase, run setup first")
		}

		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		//TODO: Run migrations down
		return nil
	},
}

func init() {
	downCmd.PersistentFlags().IntVar(&rollbackCount, "n", 1, "number of migrations to run down")
	downCmd.PersistentFlags().StringVar(&changeLogFile, "c", "./migrations/changelog.xml", "migrations changelog")
	downCmd.PersistentFlags().StringVar(&environment, "e", "development", "environment to run the migrations against")
	migrateCmd.AddCommand(downCmd)
}
