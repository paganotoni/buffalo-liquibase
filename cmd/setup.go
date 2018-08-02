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
		if runtime.GOOS == "darwin" {
			if _, err := exec.LookPath("brew"); err == nil {
				c := exec.Command("brew", "install", "liquibase")
				c.Stdin = os.Stdin
				c.Stderr = os.Stderr
				c.Stdout = os.Stdout
				return c.Run()
			}
		}

		return errors.New("Could not install liquibase cli with brew, please install manually, for instructions visit https://www.liquibase.org/")
	},
}

func init() {
	liquibaseCmd.AddCommand(setupCmd)
}
