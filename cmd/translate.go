package cmd

import (
	"github.com/paganotoni/buffalo-liquibase/liquibase"
	"github.com/spf13/cobra"
)

var fizzMigrationsPath string
var translatePath string

// translateCmd translates fizz into liquibase
var translateCmd = &cobra.Command{
	Use:     "translate",
	Short:   "description about this plugin",
	Aliases: []string{"t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		t := liquibase.NewTranslator(fizzMigrationsPath)
		return t.Translate(translatePath)
	},
}

func init() {
	translateCmd.PersistentFlags().StringVar(&fizzMigrationsPath, "path", "./migrations", "lets change the folder where fizz migrations are")
	translateCmd.PersistentFlags().StringVar(&translatePath, "output", "./migrations", "lets change the folder migrations will be outputed")
	liquibaseCmd.AddCommand(translateCmd)
}
