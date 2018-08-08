package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/gobuffalo/pop"
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

		if _, err := os.Stat(databaseYmlFile); os.IsNotExist(err) {
			return errors.New("please run this command in your buffalo app root")
		}

		if err := pop.LoadConfigFile(); err != nil {
			return err
		}

		env := pop.Connections[environment]
		if env == nil {
			return fmt.Errorf("could not find %v environment in your database.yml", environment)
		}

		originalURL := env.URL()

		r := regexp.MustCompile(`postgres:\/\/(?P<username>.*):(?P<password>.*)@(?P<host>.*):(?P<port>.*)\/(?P<database>.*)\?(?P<extras>.*)`)
		match := r.FindStringSubmatch(originalURL)
		if match == nil {
			return fmt.Errorf("could not convert %v url into liquibase", environment)
		}

		URL := fmt.Sprintf("jdbc:postgresql://%v:%v/%v?%v", match[3], match[4], match[5], match[6])

		runArgs := []string{
			"--driver=org.postgresql.Driver",
			"--url=" + URL,
			"--logLevel=debug",
			"--username=" + match[1],
			"--password=" + match[2],
			"--changeLogFile=" + changeLogFile,
			"rollbackCount",
			strconv.Itoa(rollbackCount),
		}

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
	rollbackCmd.PersistentFlags().StringVar(&databaseYmlFile, "d", "./database.yml", "database.yml file")
	liquibaseCmd.AddCommand(rollbackCmd)
}
