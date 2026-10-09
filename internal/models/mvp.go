package models

import (
	"net/http"
	"time"
	"uuid"
)

type Presenter interface {
	RegisterModel(model Model)
	DisplayHomeHandler(projectId uuid.UUID) http.Handler
	DisplayTableRow(tr TableRow) http.Handler
}

type Model interface {
	GetActivities(projectId uuid.UUID) ([]Activity, error)
	CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (Activity, error)
	CreateTableRow(projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (TableRow, error)
	GetDependencies(activityId uuid.UUID) ([]Activity, error)
	GetTableRows(projectId uuid.UUID) ([]TableRow, error)
	UpdateActivityName(activityId uuid.UUID, activityName string) (Activity, error)
	UpdateActivityDuration(activityId uuid.UUID, duration time.Duration) (Activity, error)
	UpdateActivityDependencies(activityId uuid.UUID, activityDependencyNames []string) ([]Activity, error)
}
