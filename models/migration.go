package models

import (
	"encoding/xml"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/paganotoni/buffalo-liquibase/models/liquibase"
)

type Migration struct {
	Name    string
	Version string
	Type    string

	UpSQL   template.HTML
	DownSQL template.HTML
}

func (m *Migration) XML() (string, error) {
	result, err := xml.MarshalIndent(liquibase.ChangeLog{
		ChangeSet: liquibase.ChangeSet{
			Author:  "buffalo-liquibase",
			ID:      fmt.Sprintf("%v-%v", m.Version, m.Name),
			UpSQL:   m.UpSQL,
			DownSQL: m.DownSQL,
		},

		Ns:             "http://www.liquibase.org/xml/ns/dbchangelog",
		Xsi:            "http://www.w3.org/2001/XMLSchema-instance",
		Ext:            "http://www.liquibase.org/xml/ns/dbchangelog-ext",
		SchemaLocation: `http://www.liquibase.org/xml/ns/dbchangelog http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-3.0.xsd http://www.liquibase.org/xml/ns/dbchangelog-ext http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-ext.xsd`,
	}, "  ", "    ")

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return html.UnescapeString(string(result)), nil
}

func (m Migration) Write(path string) error {
	content, err := m.XML()

	if err != nil {
		return err
	}

	if path == "./migrations" {
		if _, err := os.Stat("./migations"); os.IsNotExist(err) {
			os.Mkdir(path, 0777)
		}
	}

	return ioutil.WriteFile(fmt.Sprintf(filepath.Join(path, "%v-%v.xml"), m.Version, m.Name), []byte(content), 0644)
}
