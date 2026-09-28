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
