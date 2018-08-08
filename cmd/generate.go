package cmd

import (
	"errors"
	"path/filepath"
	"strings"
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
		nameParts := strings.Split(args[0], "/")
		migration := liquibase.Migration{
			Name:    nameParts[len(nameParts)-1],
			Version: time.Now().UTC().Format("20060102150405"),
		}

		path := append([]string{generatePath}, nameParts[:len(nameParts)-1]...)
		return migration.Write(filepath.Join(path...))
	},
}

func init() {
	generateCmd.PersistentFlags().StringVar(&generatePath, "path", "./migrations", "the path where we will generate the new migration")
	liquibaseCmd.AddCommand(generateCmd)
}
