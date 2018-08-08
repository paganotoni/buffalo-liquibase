package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "sets up liquibase and postgresql java driver",
	RunE: func(cmd *cobra.Command, args []string) error {

		if runtime.GOOS != "darwin" {
			return errors.New("Could not find liquibase, please install manually, for instructions visit https://www.liquibase.org/")
		}

		if _, err := exec.LookPath("brew"); err != nil {
			return errors.New("Could not find brew, please install liquibase manually, for instructions visit https://www.liquibase.org/")
		}

		if _, err := exec.LookPath("liquibase"); err == nil {
			cmd.Printf("Found liquibase, all good.\n")
		} else {
			c := exec.Command("brew", "install", "liquibase")
			c.Stdin = os.Stdin
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			err := c.Run()
			if err != nil {
				return err
			}
		}

		cmd.Printf("Downloading PotgreSQL driver\n")

		u, err := user.Current()
		if err != nil {
			return err
		}

		c := exec.Command("mkdir", "-p", fmt.Sprintf("/Users/%v/Library/Java/Extensions/", u.Username))
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		err = c.Run()
		if err != nil {
			return err
		}

		c = exec.Command("curl", "-s", "https://jdbc.postgresql.org/download/postgresql-42.2.4.jar", "--output", fmt.Sprintf("/Users/%v/Library/Java/Extensions/postgresql.jar", u.Username))
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
}

func init() {
	liquibaseCmd.AddCommand(setupCmd)
}
