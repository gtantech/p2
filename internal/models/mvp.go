package models

import (
	"net/http"
	"time"
	"uuid"
)

type Presenter interface {
	RegisterModel(model Model)
	HtmlHomeHandler(projectId uuid.UUID) http.Handler
	HtmlTableRow(tr TableRow) http.Handler
}

type Model interface {
	CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (Activity, error)
	CreateTableRow(projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (TableRow, error)
	GetActivities(projectId uuid.UUID) ([]Activity, error)
	GetTableRowsSortedByRank(projectId uuid.UUID) ([]TableRow, error)
	GetTableRowByActivityId(activityId uuid.UUID) (TableRow, error)
	UpdateActivityName(activityId uuid.UUID, activityName string) (Activity, error)
	UpdateActivityDuration(activityId uuid.UUID, duration time.Duration) (Activity, error)
	UpdateActivityDependencies(activityId uuid.UUID, activityDependencyNames []string) ([]Activity, error)
}
