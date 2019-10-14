package models

import (
	"encoding/xml"
)

type DatabaseChangeLog struct {
	//General XML things
	XMLName        xml.Name `xml:"databaseChangeLog"`
	Ns             string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	Ext            string   `xml:"xmlns:ext,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	//Attributes
	ChangeSet  ChangeSet    `xml:"changeSet"`
	Include    []Include    `xml:"include"`
	IncludeAll []IncludeAll `xml:"includeAll"`
}
