package models

import (
	"time"
	"uuid"
)

const (
	TableRowSortRankStep = int64(1000)
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

type PutFromRowActivityNameChange struct {
	RowId   uuid.UUID `json:"rowId"`
	DomName string    `json:"domName"`
}
