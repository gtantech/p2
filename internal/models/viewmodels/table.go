package viewmodels

import (
	"time"
	"uuid"
)

const (
	TableRowSortRankStep = int64(1000)
)

type TableRow struct {
	ActivityId            uuid.UUID
	ProjectId             uuid.UUID
	ActivityName          string
	PredecessorActivities []string
	Duration              time.Duration
	SortRank              int64
}

type Table struct {
	ProjectId uuid.UUID
	Rows      []TableRow
}

type InputAutocomplete string

const (
	On  InputAutocomplete = "on"
	Off InputAutocomplete = "off"
)
