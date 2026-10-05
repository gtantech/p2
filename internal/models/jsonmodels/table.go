package jsonmodels

import (
	"uuid"
)

type CreateEmptyTableRow struct {
	PreviousRowSortRank *int64    `json:"previousRowSortRank"`
	NextRowSortRank     *int64    `json:"nextRowSortRank"`
	ProjectId           uuid.UUID `json:"projectId"`
}

type UpdateActivityFromInput struct {
	ActivityId uuid.UUID `json:"activityId"`
	DomName    string    `json:"domName"`
}

type UpdateDependencyFromInput struct {
	ActivityId uuid.UUID `json:"activityId"`
	ProjectId  uuid.UUID `json:"projectId"`
	DomName    string    `json:"domName"`
}

type AddActivitiesFromRow struct {
	IntoProjectId      uuid.UUID `json:"intoProjectId"`
	FromRowActivityId  uuid.UUID `json:"fromRowActivityId"`
	ActivityNamesToAdd []string  `json:"activityNamesToAdd"`
}
