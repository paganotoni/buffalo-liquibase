package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:                "test",
	Short:              "test runs migrations up and then runs tests",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		conn := pop.Connections["test"]
		err := pop.DropDB(conn)
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			return err
		}

		err = pop.CreateDB(conn)
		if err != nil {
			return err
		}

		environment = "test"
		up, err := buildUpCommand()
		if err != nil {
			return err
		}

		err = up.Run()
		if err != nil {
			return err
		}

		os.Setenv("GO_ENV", "test")

		if len(args) == 0 {
			args = []string{"./..."}
		}

		runArgs := append([]string{"test", "-p", "1"}, args...)
		c := exec.Command("go", runArgs...)
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout

		fmt.Println(strings.Join(append([]string{"go"}, runArgs...), " "))

		return c.Run()
	},
}

func init() {
	liquibaseCmd.AddCommand(testCmd)
}
