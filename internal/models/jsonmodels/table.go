package jsonmodels

import (
	"time"
	"uuid"
)

type CreateEmptyTableRow struct {
	PreviousRowSortRank *int64    `json:"previousRowSortRank"`
	NextRowSortRank     *int64    `json:"nextRowSortRank"`
	ProjectId           uuid.UUID `json:"projectId"`
}

type UpdateActivity struct {
	ActivityId  uuid.UUID     `json:"activityId"`
	DisplayName string        `json:"displayName"`
	Duration    time.Duration `json:"duration"`
}
