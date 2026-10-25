package models

import (
	"time"
	"uuid"
)

type Table struct {
	Rows []TableRow
}

type TableRow struct {
	RowId           uuid.UUID
	ActivityName    string
	DependencyNames []string
	Duration        time.Duration
}
