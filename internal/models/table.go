package models

import "uuid"

type Table struct {
	ProjectId uuid.UUID
	Rows      []TableRow
}

type TableRow struct {
	Activity
	ProjectId    uuid.UUID
	Dependencies []string
	SortRank     int64
}
