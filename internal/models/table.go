package models

import (
	"time"
	"uuid"
)

type Table struct {
	Rows []TableRow
}

type TableRow struct {
	RowId        uuid.UUID
	Activity     string
	Dependencies []string
	Duration     time.Duration
}
