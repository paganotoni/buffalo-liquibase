package migrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	r := require.New(t)
	path := filepath.Join(os.TempDir(), "changelog.xml")

	f, err := os.Create(path)
	r.NoError(err)

	baseXML := `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<databaseChangeLog xmlns="http://www.liquibase.org/xml/ns/dbchangelog" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-2.0.xsd">
		<include file="migrations/schema/20190625162047-add_uuid_extension.xml" />
		<include file="migrations/schema/20190625162553-create_devices.xml" />
	</databaseChangeLog>
	`
	_, err = f.Write([]byte(baseXML))
	r.NoError(err)

	chl, err := Parser{}.ParseXML(path)
	r.NoError(err)
	r.Len(chl.Include, 2)

	r.Equal("migrations/schema/20190625162047-add_uuid_extension.xml", chl.Include[0].File)
	r.Equal("migrations/schema/20190625162553-create_devices.xml", chl.Include[1].File)
}

func TestParserChangeSet(t *testing.T) {
	r := require.New(t)
	path := filepath.Join(os.TempDir(), "changeset.xml")

	f, err := os.Create(path)
	r.NoError(err)

	baseXML := `
	<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<databaseChangeLog xmlns="http://www.liquibase.org/xml/ns/dbchangelog" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-2.0.xsd">
		<changeSet id="20191010215845-update_residual_values" author="buffalo-liquibase">
			<sql>
				UPDATE replacement_periods SET residual_value = 0.5210 WHERE length = 12;
				UPDATE replacement_periods SET residual_value = 0.7030 WHERE length = 24;
				UPDATE replacement_periods SET residual_value = 0.8010 WHERE length = 36;
			</sql>
			<rollback></rollback>
		</changeSet>
	</databaseChangeLog>
	`
	_, err = f.Write([]byte(baseXML))
	r.NoError(err)

	chl, err := Parser{}.ParseXML(path)
	r.NoError(err)
	r.Equal("20191010215845-update_residual_values", chl.ChangeSet.ID)
	r.Equal("", string(chl.ChangeSet.DownSQL))
	r.Contains(chl.ChangeSet.UpSQL, "UPDATE replacement_periods SET residual_value = 0.5210 WHERE length = 12;")

}
