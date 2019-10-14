package migrator

import (
	"os"

	"github.com/gobuffalo/pop"
	_ "github.com/lib/pq"
	"github.com/paganotoni/buffalo-liquibase/liquibase/models"
	"github.com/pkg/errors"

	"github.com/paganotoni/buffalo-liquibase/liquibase/models"
)

// PostgresMigrator interface check
var _ Migrator = (*PostgresMigrator)(nil)

const (
	migrationTablesQuery = `SELECT 'public.databasechangelog'::regclass, 'public.databasechangeloglock'::regclass;`
	migrationsTablesDDL  = `
	CREATE TABLE public.databasechangelog (
		id character varying(255) NOT NULL,
		author character varying(255) NOT NULL,
		filename character varying(255) NOT NULL,
		dateexecuted timestamp without time zone NOT NULL,
		orderexecuted integer NOT NULL,
		exectype character varying(10) NOT NULL,
		md5sum character varying(35),
		description character varying(255),
		comments character varying(255),
		tag character varying(255),
		liquibase character varying(20),
		contexts character varying(255),
		labels character varying(255),
		deployment_id character varying(10)
	);
	
	CREATE TABLE public.databasechangeloglock (
		id integer NOT NULL,
		locked boolean NOT NULL,
		lockgranted timestamp without time zone,
		lockedby character varying(255)
	);
	`

	lockStatement = "INSERT INTO databasechangeloglock (id, locked, lockedby) VALUES (?, ?, ?)"

	migrationsLogs = "SELECT * FROM databasechangelog"
)

// PostgresMigrator is the implementation of the migrator for the postgres
// database type. At the type i'm writing this there is no need for other
// kind of migrators but i preffer to leave the door open.
type PostgresMigrator struct {
	Conn *pop.Connection

	log models.DatabaseChangeLog
}

func (pm PostgresMigrator) ensureTables() error {
	//Check table existence
	err := pm.Conn.RawQuery(migrationTablesQuery).Exec()
	if err == nil {
		return nil
	}

	//Create tables
	err = pm.Conn.RawQuery(migrationsTablesDDL).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (pm PostgresMigrator) canMigrate() (bool, error) {
	result := struct {
		Count int64 `db:"count"`
	}{}

	err := pm.Conn.RawQuery("SELECT count(locked = true) FROM databasechangeloglock;").First(&result)
	return result.Count == 0, err
}

func (pm PostgresMigrator) lock() error {
	return pm.Conn.Transaction(func(tx *pop.Connection) error {
		return tx.RawQuery(lockStatement, os.Getpid(), true, "buffalo-liquibase").Exec()
	})
}

func (pm PostgresMigrator) getMigrationLogs() (models.MigrationLogs, error) {
	mlogs := models.MigrationLogs{}
	err := pm.Conn.Transaction(func(tx *pop.Connection) error {
		return tx.All(&mlogs)
	})

	return mlogs, err
}

func (pm PostgresMigrator) Prepare() error {
	err := pm.ensureTables()
	if err != nil {
		return err
	}

	return nil
}

func (pm PostgresMigrator) Up() error {
	//1. Run Prepare
	if err := pm.Prepare(); err != nil {
		return err
	}

	//2. Ensure we can proceed (no lock)
	proceed, err := pm.canMigrate()
	if err != nil {
		return err
	}

	if !proceed {
		return errors.New("Database is locked by another Liquibase process at this time, migrations cannot run")
	}

	//3. Acquire lock in table (write it)
	err = pm.lock()
	if err != nil {
		return errors.Wrap(err, "error while locking")
	}

	//4. Parse migrations
	//5. Find last migration
	//6. Run missing and add to migrations registry
	//7. Release lock
	return nil
}

func (pm PostgresMigrator) Down(count int) error {
	//Run migrations down
	return nil
}
