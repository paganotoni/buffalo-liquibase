package models

import "html/template"

type ChangeSet struct {
	ID      string        `xml:"id,attr"`
	Author  string        `xml:"author,attr"`
	UpSQL   template.HTML `xml:"sql"`
	DownSQL template.HTML `xml:"rollback"`
}
