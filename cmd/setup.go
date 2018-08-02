package cmd

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// setupCmd translates fizz into liquibase
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "description about this plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("liquibase"); err == nil {
			cmd.Printf("Found liquibase, all good.")
		}

		if runtime.GOOS != "darwin" {
			return errors.New("Could not find liquibase, please install manually, for instructions visit https://www.liquibase.org/")
		}

		if _, err := exec.LookPath("brew"); err != nil {
			return errors.New("Could not find liquibase or brew, please install manually, for instructions visit https://www.liquibase.org/")
		}

		c := exec.Command("brew", "install", "liquibase")
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()

	},
}

func init() {
	liquibaseCmd.AddCommand(setupCmd)
}
