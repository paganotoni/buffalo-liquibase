package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/gobuffalo/pop"
	"github.com/paganotoni/buffalo-liquibase/liquibase"
	"github.com/spf13/cobra"
)

var changeLogFile string
var environment string
var databaseYmlFile string

// upCmd runs /migrations or --path up against buffalo db with liquibase
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "runs liquibase migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("liquibase"); err != nil {
			return errors.New("could not find liquibase, run setup first")
		}

		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		c, err := buildUpCommand()
		if err != nil {
			return err
		}

		return c.Run()
	},
}

func buildUpCommand() (*exec.Cmd, error) {
	runArgs, err := liquibase.BuildRunArgsFor(environment)
	if err != nil {
		return nil, err
	}

	runArgs = append(runArgs, []string{
		"--changeLogFile=" + changeLogFile,
		"update",
	}...)

	c := exec.Command("liquibase", runArgs...)
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	return c, nil
}

func init() {
	upCmd.PersistentFlags().StringVar(&changeLogFile, "c", "./migrations/changelog.xml", "migrations changelog")
	upCmd.PersistentFlags().StringVar(&environment, "e", "development", "environment to run the migrations against")
	liquibaseCmd.AddCommand(upCmd)
}
