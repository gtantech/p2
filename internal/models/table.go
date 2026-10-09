package models

import "uuid"

const (
	TableRowSortRankStep = int64(1000)
)

type Table struct {
	ProjectId uuid.UUID
	Rows      []TableRow
}

type TableRow struct {
	*Activity
	ProjectId    uuid.UUID
	Dependencies []string
	SortRank     int64
}
