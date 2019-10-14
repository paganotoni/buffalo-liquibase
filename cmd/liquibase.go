package cmd

import (
	"os"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/defaults"
	"github.com/spf13/cobra"
)

var all bool
var env string

// liquibaseCmd represents the liquibase command
var liquibaseCmd = &cobra.Command{
	Use:     "liquibase",
	Short:   "Runs yours migrations using liquibase",
	Aliases: []string{"db"},
	PersistentPreRun: func(c *cobra.Command, args []string) {
		if !c.PersistentFlags().Changed("env") {
			env = defaults.String(os.Getenv("GO_ENV"), env)
		}

		pop.LoadConfigFile()
	},
}

func init() {
	liquibaseCmd.PersistentFlags().StringVarP(&env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")
	rootCmd.AddCommand(liquibaseCmd)
}
