package models

import (
	"time"
	"uuid"
)

type Table struct {
	Rows []TableRows
}

type TableRows struct {
	RowId        uuid.UUID
	Activity     string
	Dependencies []string
	Duration     time.Duration
}
