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
	// CreateActivity returns ErrProjectNotFound if projectId does not exist.
	CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (Activity, error)

	// CreateTableRow returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId does not exist.
	CreateTableRow(projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (TableRow, error)

	// GetActivities returns ErrProjectNotFound if projectId does not exist.
	GetActivities(projectId uuid.UUID) ([]Activity, error)

	// GetTableRowsSortedByRank returns ErrProjectNotFound if projectId does not exist.
	GetTableRowsSortedByRank(projectId uuid.UUID) ([]TableRow, error)

	// GetTableRowByActivityId returns ErrActivityNotFound if activityId does not exist.
	GetTableRowByActivityId(activityId uuid.UUID) (TableRow, error)

	// GetNextTableRowByActivityId returns ErrEndOfTable if the row specified by activityId is the last row,
	// ErrProjectNotFound if projectId does not exist, and ErrActivityNotFound if activityId does not exist.
	GetNextTableRowByActivityId(projectId uuid.UUID, activityId uuid.UUID) (TableRow, error)

	// UpdateActivityName returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityName(activityId uuid.UUID, activityName string) (Activity, error)

	// UpdateActivityDuration returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityDuration(activityId uuid.UUID, duration time.Duration) (Activity, error)

	// UpdateActivityDependencies returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityDependencies(activityId uuid.UUID, activityDependencyNames []string) ([]Activity, error)
}
