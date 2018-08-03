package liquibase

import (
	"encoding/xml"
	"html/template"
)

type ChangeSet struct {
	ID      string        `xml:"id,attr"`
	Author  string        `xml:"author,attr"`
	UpSQL   template.HTML `xml:"sql"`
	DownSQL template.HTML `xml:"rollback"`
}

type DatabaseChangeLog struct {
	ChangeSet ChangeSet `xml:"changeSet"`

	XMLName        xml.Name `xml:"databaseChangeLog"`
	Ns             string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	Ext            string   `xml:"xmlns:ext,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
}
