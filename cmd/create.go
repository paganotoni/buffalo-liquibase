package cmd

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

// createCmd creates your databases
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates databases for you",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if !all {
			return pop.CreateDB(getConn())
		}

		for _, conn := range pop.Connections {
			err = pop.CreateDB(conn)
			if err != nil {
				return err
			}
		}

		return nil

	},
}

func init() {
	createCmd.Flags().BoolVarP(&all, "all", "a", false, "Creates all of the databases in the database.yml")
	liquibaseCmd.AddCommand(createCmd)
}

func getConn() *pop.Connection {
	conn := pop.Connections[env]
	if conn == nil {
		fmt.Printf("There is no connection named %s defined!\n", env)
		os.Exit(1)
	}
	return conn
}
