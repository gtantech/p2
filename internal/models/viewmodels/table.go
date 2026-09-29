package viewmodels

import (
	"time"
	"uuid"
)

type TableRow struct {
	ActivityId            uuid.UUID
	ActivityName          string
	PredecessorActivities []string
	Duration              time.Duration
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
