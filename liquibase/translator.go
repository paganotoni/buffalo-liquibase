package liquibase

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

//Translator takes the job of translating fizz migrations
type Translator struct {
	path       string
	migrations []pop.Migration
}

func NewTranslator(path string) Translator {
	t := Translator{
		path: path,
	}

	t.LoadMigrations()
	return t
}

func (t *Translator) LoadMigrations() {
	dir := t.path
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return
	}

	var mrx = regexp.MustCompile(`^(\d+)_([^\.]+)(\.[a-z0-9]+)?\.(up|down)\.(sql|fizz)$`)

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches := mrx.FindAllStringSubmatch(info.Name(), -1)
			if len(matches) == 0 {
				return nil
			}

			m := matches[0]
			mf := pop.Migration{
				Path:      p,
				Version:   m[1],
				Name:      m[2],
				Direction: m[4],
				Type:      m[5],
			}

			t.migrations = append(t.migrations, mf)
		}

		return nil
	})
}

func (t *Translator) Translate() error {
	fmt.Printf("Found %v migrations in %v\n", len(t.migrations), t.path)

	for _, mi := range t.migrations {
		err := t.translateMigration(mi, t.findDownFor(mi))
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Translator) findDownFor(mi pop.Migration) pop.Migration {
	for _, dmi := range t.migrations {
		if dmi.Direction == "up" {
			continue
		}

		if mi.Name == dmi.Name {
			return dmi
		}
	}
	return pop.Migration{}
}

func (t *Translator) translateMigration(up, down pop.Migration) error {
	upsql, err := t.convertMigration(up)
	if err != nil {
		return err
	}

	downsql, err := t.convertMigration(down)
	if err != nil {
		return err
	}

	return t.renderFile(up.Name, up.Version, upsql, downsql)
}

func (t *Translator) convertMigration(mi pop.Migration) (string, error) {
	f, err := os.Open(mi.Path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	up := string(b)
	if mi.Type == "fizz" {
		up, err = fizz.AString(up, translators.NewPostgres())
		if err != nil {
			return "", errors.Wrapf(err, "could not fizz the migration %s", mi.Path)
		}
	}

	return up, nil
}

func (t *Translator) renderFile(name, version, up, down string) error {
	midata := struct {
		Name      string
		Timestamp string
		UpSQL     template.HTML
		DownSQL   template.HTML
	}{name, version, template.HTML(up), template.HTML(down)}

	tmp, _ := template.New("xml").Parse(migrationXMLTemplate)
	var tpl bytes.Buffer
	if err := tmp.Execute(&tpl, midata); err != nil {
		return err
	}

	result := tpl.String()
	result = xmlfmt.FormatXML(result, "\t", "  ")

	if _, err := os.Stat("./migrations"); os.IsNotExist(err) {
		os.Mkdir("./migrations", 0777)
	}

	data := []byte(result)
	return ioutil.WriteFile(fmt.Sprintf("./migrations/%v-%v.xml", midata.Timestamp, midata.Name), data, 0644)
}

const migrationXMLTemplate = `
<databaseChangeLog xmlns="http://www.liquibase.org/xml/ns/dbchangelog" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:ext="http://www.liquibase.org/xml/ns/dbchangelog-ext" xsi:schemaLocation="http://www.liquibase.org/xml/ns/dbchangelog http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-3.0.xsd http://www.liquibase.org/xml/ns/dbchangelog-ext http://www.liquibase.org/xml/ns/dbchangelog/dbchangelog-ext.xsd">
	<changeSet author="buffalo-liquibase" id="{{.Timestamp}}-{{.Name}}">
        <sql>
{{.UpSQL}}
		</sql>
		<rollback>
{{.DownSQL}}
	    </rollback>
    </changeSet>
</databaseChangeLog>
`
