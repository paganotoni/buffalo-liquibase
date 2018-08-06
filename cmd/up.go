package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var changeLogFile string

// upCmd runs /migrations or --path up against buffalo db with liquibase
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "runs liquibase migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("liquibase"); err != nil {
			return errors.New("could not find liquibase, run setup first")
		}

		if _, err := os.Stat("database.yml"); os.IsNotExist(err) {
			return errors.New("please run this command in your buffalo app root")
		}

		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		originalURL := pop.Connections["development"].URL()
		r := regexp.MustCompile(`postgres:\/\/(?P<username>.*):(?P<password>.*)@(?P<host>.*):(?P<port>.*)\/(?P<database>.*)\?.*`)
		match := r.FindStringSubmatch(originalURL)
		URL := fmt.Sprintf("jdbc:postgresql://%v:%v/%v?ssl=true&sslfactory=org.postgresql.ssl.NonValidatingFactory", match[3], match[4], match[5])

		runArgs := []string{
			"--driver=org.postgresql.Driver",
			"--url=" + URL,
			"--logLevel=debug",
			"--username=" + match[1],
			"--password=" + match[2],
			"--changeLogFile=" + changeLogFile,
			"update",
		}

		c := exec.Command("liquibase", runArgs...)
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		return c.Run()
	},
}

func init() {
	translateCmd.PersistentFlags().StringVar(&changeLogFile, "c", "./migrations/changelog.xml", "migrations changelog")
	liquibaseCmd.AddCommand(upCmd)
}
