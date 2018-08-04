package cmd

import (
	"errors"
	"time"

	"github.com/paganotoni/buffalo-liquibase/liquibase"

	"github.com/spf13/cobra"
)

var generatePath string

//generate a sample liquibase migration in /migrations or --path
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates a new liquibase xml migration",
	Aliases: []string{"g"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must provide a name for the migration")
		}

		migration := liquibase.Migration{
			Name:    args[0],
			Version: time.Now().UTC().Format("20060102150405"),
		}

		return migration.Write(generatePath)
	},
}

func init() {
	generateCmd.PersistentFlags().StringVar(&generatePath, "path", "./migrations", "the path where we will generate the new migration")
	liquibaseCmd.AddCommand(generateCmd)
}
