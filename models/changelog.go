package models

import (
	"time"

	"github.com/gobuffalo/nulls"
)

// ChangeLog ...
type ChangeLog struct {
	ID            string       `db:"id"`
	Author        string       `db:"author"`
	Filename      string       `db:"filename"`
	DateExecuted  time.Time    `db:"dateexecuted"`
	OrderExecuted int          `db:"orderexecuted"`
	ExecType      string       `db:"exectype"`
	MD5Sum        nulls.String `db:"md5sum"`
	Description   nulls.String `db:"description"`
	Comments      nulls.String `db:"comments"`
	Tag           nulls.String `db:"tag"`
	Liquibase     nulls.String `db:"liquibase"`
	Contexts      nulls.String `db:"contexts"`
	Labels        nulls.String `db:"labels"`
	DeploymentID  nulls.String `db:"deployment_id"`
}

// MigrationLogs ...
type ChangeLogs []ChangeLog

// TableName function ...
func (ml ChangeLog) TableName() string {
	return "databasechangelog"
}
