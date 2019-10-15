package migrator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/suite"
	"github.com/stretchr/testify/require"

	"github.com/paganotoni/buffalo-liquibase/models"
)

type PostgresSuite struct {
	*suite.Model

	Migrator *PostgresMigrator
}

func TestPostgres(t *testing.T) {

	conn, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect: "postgres",
		URL:     "postgres://postgres:postgres@127.0.0.1:5432/buffalo-liquibase?sslmode=disable",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir("migrations", 0777)
	if err != nil {
		t.Fail()
		return
	}
	defer os.RemoveAll("migrations")

	pop.Connections["test"] = conn

	model := suite.NewModel()
	model.DB = conn
	model.Assertions = require.New(t)

	ps := &PostgresSuite{
		Model: model,

		Migrator: &PostgresMigrator{
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
	ps.NoError(ps.Migrator.Prepare())
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

func (ps PostgresSuite) Test_Unlock() {
	ps.NoError(ps.Migrator.Prepare())

	err := ps.DB.RawQuery("INSERT INTO databasechangeloglock (id, locked) values (1, true);").Exec()
	ps.NoError(err)

	err = ps.Migrator.unlock()
	ps.NoError(err)

	result := struct {
		Count int
	}{}

	err = ps.DB.RawQuery("SELECT count(*) FROM databasechangeloglock;").First(&result)
	ps.NoError(err)
	ps.Equal(0, result.Count)
}

func (ps PostgresSuite) Test_GetMigrationLogs() {
	log := models.ChangeLog{
		ID:            "20190625162047-add_uuid_extension",
		Author:        "buffalo-liquibase",
		Filename:      "/migrations/schema/20190625162047-add_uuid_extension.xml",
		DateExecuted:  time.Now(),
		OrderExecuted: 1,
		ExecType:      "EXECUTED",
		MD5Sum:        nulls.NewString("8:6b44712359cb1cea8882505ea4ce8649"),
		Description:   nulls.NewString("sql"),
		Comments:      nulls.NewString(""),
		Tag:           nulls.NewString(""),
		Liquibase:     nulls.NewString("3.8.0"),
		Contexts:      nulls.NewString(""),
		Labels:        nulls.NewString(""),
		DeploymentID:  nulls.NewString("9354664289"),
	}
	ps.DB.Create(&log)

	errLogs := ps.Migrator.loadDatabaseChangelog()
	ps.NoError(errLogs)

	ps.Len(ps.Migrator.databaseChangelog, 1)

	ps.Equal("20190625162047-add_uuid_extension", ps.Migrator.databaseChangelog[0].ID)
	ps.Equal("buffalo-liquibase", ps.Migrator.databaseChangelog[0].Author)
}

func (ps PostgresSuite) Test_PendingMigrations() {
	path := filepath.Join("migrations", "changelog.xml")

	f, err := os.Create(path)
	ps.NoError(err)

	baseXML := `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<databaseChangeLog xmlns="http://www.liquibase.org/xml/ns/dbchangelog" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-2.0.xsd">
		<include file="migrations/schema/20190625162047-add_uuid_extension.xml" />
		<include file="migrations/schema/20190625162553-create_devices.xml" />
	</databaseChangeLog>
	`
	_, err = f.Write([]byte(baseXML))
	ps.NoError(err)

	ps.NoError(ps.Migrator.Prepare())
	err = ps.Migrator.loadDatabaseChangelog()
	ps.NoError(err)

	ps.Migrator.changelog, err = Parser{}.ParseXML("migrations/changelog.xml")
	ps.NoError(err)

	chl := ps.Migrator.pendingMigrations()
	ps.Len(chl, 2)

	ps.Migrator.databaseChangelog = append(ps.Migrator.databaseChangelog, models.ChangeLog{
		Filename: "migrations/schema/20190625162553-create_devices.xml",
	})

	chl = ps.Migrator.pendingMigrations()
	ps.Len(chl, 1)

	ps.Migrator.databaseChangelog = models.ChangeLogs{
		{Filename: "migrations/schema/20190625162553-create_devices_non_in_the_list.xml"},
	}

	chl = ps.Migrator.pendingMigrations()
	ps.Len(chl, 2)

}
