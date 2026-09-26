package models

import (
	"time"
	"uuid"
)

type ViewActivity struct {
	Id          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type ViewTable struct {
	Rows []*ViewTableRow
}

type ViewTableRow struct {
	Activity     *ViewActivity
	Dependencies []*ViewActivity
}

type ViewDisplayEmptyTableRowParams struct {
	ActivityId uuid.UUID
	ProjectId  uuid.UUID
}
