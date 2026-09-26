package models

import (
	"time"
	"uuid"
)

type Activity struct {
	Id          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type Table struct {
	Rows []*TableRow
}

type TableRow struct {
	Activity     *Activity
	Dependencies []*Activity
}
