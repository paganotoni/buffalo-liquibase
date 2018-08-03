package liquibase

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/pkg/errors"
)

//Translator takes the job of translating fizz migrations
type Translator struct {
	path       string
	migrations []Migration
}

func NewTranslator(path string) Translator {
	t := Translator{
		path: path,
	}

	t.LoadMigrations()
	return t
}

func (t *Translator) Translate() error {
	fmt.Printf("Found %v migrations in %v\n", len(t.migrations), t.path)

	for _, mi := range t.migrations {
		err := mi.Write()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Translator) Convert(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(path, ".fizz") {
		return string(b), nil
	}

	sql, err := fizz.AString(string(b), translators.NewPostgres())
	if err != nil {
		return "", errors.Wrapf(err, "could not fizz the migration %s", path)
	}

	return sql, nil
}

func (t *Translator) LoadMigrations() {
	dir := t.path
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return
	}

	var mrx = regexp.MustCompile(`^(\d+)_([^\.]+)(\.[a-z0-9]+)?\.(up|down)\.(sql|fizz)$`)

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		matches := mrx.FindAllStringSubmatch(info.Name(), -1)
		if len(matches) == 0 {
			return nil
		}

		sql, err := t.Convert(p)
		esql := template.HTML(sql)

		m := matches[0]
		for index, mi := range t.migrations {
			if mi.Name != m[2] || mi.Version != m[1] {
				continue
			}

			if m[4] == "down" {
				t.migrations[index].DownSQL = esql
			} else {
				t.migrations[index].UpSQL = esql
			}

			return nil
		}

		mf := Migration{
			Version: m[1],
			Name:    m[2],
			Type:    m[5],
		}

		if m[4] == "down" {
			mf.DownSQL = esql
		} else {
			mf.UpSQL = esql
		}

		t.migrations = append(t.migrations, mf)
		return nil
	})
}
