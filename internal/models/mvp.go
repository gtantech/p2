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
	HtmxPostActivityDependencyHandlerResp(suggestedActivityNames []string, rowActivityId uuid.UUID, projectId uuid.UUID) http.Handler
	HtmxPostActivityDependencySuggestionAddNewHandlerResp(createdTableRows []TableRow, rowActivityId uuid.UUID) http.Handler
}

type Model interface {
	// CreateActivity returns ErrProjectNotFound if projectId does not exist.
	CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (Activity, error)

	// CreateTableRow returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId does not exist.
	CreateTableRow(projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (TableRow, error)

	// CreateTableRowAtBottom returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId does not exist.
	CreateTableRowAtBottom(projectId uuid.UUID, activityId uuid.UUID) (TableRow, error)

	// CreateTableRowBelow returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId or belowActivityId does not exist.
	CreateTableRowBelow(projectId uuid.UUID, activityId uuid.UUID, belowActivityId uuid.UUID) (TableRow, error)

	// CreateTableRowBelow returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if successorActivityId or predecessorActivityId does not exist.
	CreateDependency(projectId uuid.UUID, successorActivityId uuid.UUID, predecessorActivityId uuid.UUID) (Dependency, error)

	// GetActivities returns ErrProjectNotFound if projectId does not exist.
	GetActivities(projectId uuid.UUID) ([]Activity, error)

	// GetTableRowsSortedByRank returns ErrProjectNotFound if projectId does not exist.
	GetTableRowsSortedByRank(projectId uuid.UUID) ([]TableRow, error)

	// UpdateActivityName returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityName(activityId uuid.UUID, activityName string) (Activity, error)

	// UpdateActivityDuration returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityDuration(activityId uuid.UUID, duration time.Duration) (Activity, error)

	// UpdateActivityDependencies returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityDependencies(activityId uuid.UUID, dependencyActivityIds []uuid.UUID) ([]Activity, error)
}
