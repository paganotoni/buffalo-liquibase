package migrator

import (
	"os"

	"github.com/gobuffalo/pop"
	_ "github.com/lib/pq"
	"github.com/paganotoni/buffalo-liquibase/models"
	"github.com/paganotoni/buffalo-liquibase/models/liquibase"
	"github.com/pkg/errors"
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
)

// PostgresMigrator is the implementation of the migrator for the postgres
// database type. At the type i'm writing this there is no need for other
// kind of migrators but i preffer to leave the door open.
type PostgresMigrator struct {
	Conn *pop.Connection

	changelogPath string
	changelog     liquibase.ChangeLog

	//Holds the list of migrations that have already run from the database.
	databaseChangelog models.ChangeLogs
}

func (pm *PostgresMigrator) ensureTables() error {
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

func (pm *PostgresMigrator) canMigrate() (bool, error) {
	result := struct {
		Count int64 `db:"count"`
	}{}

	err := pm.Conn.RawQuery("SELECT count(locked = true) FROM databasechangeloglock;").First(&result)
	return result.Count == 0, err
}

func (pm *PostgresMigrator) lock() error {
	return pm.Conn.Transaction(func(tx *pop.Connection) error {
		return tx.RawQuery(lockStatement, os.Getpid(), true, "buffalo-liquibase").Exec()
	})
}

func (pm *PostgresMigrator) unlock() error {
	return pm.Conn.Transaction(func(tx *pop.Connection) error {
		return tx.RawQuery("DELETE FROM databasechangeloglock;").Exec()
	})
}

func (pm *PostgresMigrator) loadDatabaseChangelog() error {
	mlogs := models.ChangeLogs{}

	err := pm.Conn.All(&mlogs)
	pm.databaseChangelog = mlogs

	return err
}

func (pm *PostgresMigrator) Prepare() error {
	err := pm.ensureTables()
	if err != nil {
		return err
	}

	return nil
}

func (pm *PostgresMigrator) Up() error {

	if err := pm.Prepare(); err != nil {
		return err
	}

	proceed, err := pm.canMigrate()
	if err != nil {
		return err
	}

	if !proceed {
		return errors.New("Database is locked by another Liquibase process at this time, migrations cannot run")
	}

	pm.changelog, err = Parser{}.ParseXML(pm.changelogPath)
	if err != nil {
		return errors.Wrap(err, "error reading changelog")
	}

	err = pm.loadDatabaseChangelog()
	if err != nil {
		return err
	}

	//Get missing migrations
	//Find last order executed number

	// err = pm.lock()
	// if err != nil {
	// 	return errors.Wrap(err, "error while locking")
	// }

	// //Unlocking tables
	// pm.unlock()
	return nil
}

func (pm *PostgresMigrator) pendingMigrations() models.ChangeLogs {
	pending := models.ChangeLogs{}

xml:
	for _, cl := range pm.changelog.Include {
		for _, db := range pm.databaseChangelog {
			if cl.File == db.Filename {
				continue xml
			}
		}

		pending = append(pending, models.ChangeLog{
			Filename: cl.File,
		})
	}

	return pending
}

func (pm *PostgresMigrator) Down(count int) error {
	//Run migrations down
	return nil
}