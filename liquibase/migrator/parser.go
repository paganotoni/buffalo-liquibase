package migrator

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/paganotoni/buffalo-liquibase/liquibase/models"
)

type Parser struct{}

func (p Parser) ParseXML(path string) (models.DatabaseChangeLog, error) {

	result := models.DatabaseChangeLog{}

	file, err := os.Open(path)
	if err != nil {
		return result, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return result, err
	}

	err = xml.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
