package cmd

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	"github.com/gobuffalo/pop"
	"github.com/paganotoni/buffalo-liquibase/liquibase"
	"github.com/spf13/cobra"
)

var rollbackCount int

// rollbackCmd runs /migrations down
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollbacks migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("liquibase"); err != nil {
			return errors.New("could not find liquibase, run setup first")
		}

		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		runArgs, err := liquibase.BuildRunArgsFor(environment)
		if err != nil {
			return err
		}

		runArgs = append(runArgs, []string{
			"--changeLogFile=" + changeLogFile,
			"rollbackCount",
			strconv.Itoa(rollbackCount),
		}...)

		cmd.Println(runArgs)

		c := exec.Command("liquibase", runArgs...)
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
}

func init() {
	rollbackCmd.PersistentFlags().IntVar(&rollbackCount, "n", 1, "number of migrations to run down")
	rollbackCmd.PersistentFlags().StringVar(&changeLogFile, "c", "./migrations/changelog.xml", "migrations changelog")
	rollbackCmd.PersistentFlags().StringVar(&environment, "e", "development", "environment to run the migrations against")
	liquibaseCmd.AddCommand(rollbackCmd)
}
