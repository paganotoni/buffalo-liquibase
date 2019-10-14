package migrator

import (
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/suite"
	"github.com/stretchr/testify/require"
)

type PostgresSuite struct {
	*suite.Model

	Migrator PostgresMigrator
}

func TestPostgres(t *testing.T) {

	conn, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect: "postgres",
		URL:     "postgres://postgres:postgres@127.0.0.1:5432/buffalo-liquibase?sslmode=disable",
	})

	if err != nil {
		t.Fatal(err)
	}

	pop.Connections["test"] = conn

	model := suite.NewModel()
	model.DB = conn
	model.Assertions = require.New(t)

	ps := &PostgresSuite{
		Model: model,

		Migrator: PostgresMigrator{
			Conn: conn,
		},
	}

	suite.Run(t, ps)
}

func (ps *PostgresSuite) SetupTest() {
	if ps.DB == nil {
		return
	}

	pop.CreateDB(ps.DB)
	err := ps.DB.TruncateAll()
	ps.NoError(err)

	ps.Migrator = PostgresMigrator{
		Conn: ps.DB,
	}
}

func (ps PostgresSuite) Test_Prepare() {
	ps.NoError(ps.Migrator.Prepare())
	ps.NoError(ps.DB.RawQuery(migrationTablesQuery).Exec())

	ps.DB.RawQuery("DROP TABLE databasechangeloglock;").Exec()
	ps.DB.RawQuery("DROP TABLE databasechangelog").Exec()

	ps.NoError(ps.Migrator.Prepare())
	ps.NoError(ps.DB.RawQuery(migrationTablesQuery).Exec())
}

func (ps PostgresSuite) Test_CanMigrate() {
	result, err := ps.Migrator.canMigrate()
	ps.NoError(err)
	ps.True(result)

	err = ps.DB.RawQuery("INSERT INTO databasechangeloglock (id, locked) values (1, true);").Exec()
	ps.NoError(err)

	result, err = ps.Migrator.canMigrate()
	ps.NoError(err)
	ps.False(result)
}

func (ps PostgresSuite) Test_Lock() {
	ps.NoError(ps.Migrator.Prepare())
	ps.NoError(ps.Migrator.lock())

	result := struct {
		LockedBy    string
		Locked      bool
		LockGranted nulls.Time `db:"lockgranted"`
		ID          int
	}{}

	err := ps.DB.RawQuery("SELECT * FROM databasechangeloglock Limit 1;").First(&result)
	ps.NoError(err)
	ps.True(result.Locked)
	ps.Equal("buffalo-liquibase", result.LockedBy)
}
