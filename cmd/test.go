package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/envy"
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
		//TODO: run migrations up

		os.Setenv("GO_ENV", "test")

		if len(args) == 0 {
			args, err = findTestPackages()

			if err != nil {
				return err
			}
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

func findTestPackages() ([]string, error) {
	args := []string{}
	out, err := exec.Command(envy.Get("GO_BIN", "go"), "list", "./...").Output()
	if err != nil {
		return args, err
	}

	pkgs := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	for _, p := range pkgs {
		if !strings.Contains(string(p), "/vendor/") {
			args = append(args, string(p))
		}
	}

	return args, err
}
