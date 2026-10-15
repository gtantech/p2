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

	//Creates a project and returns the newly created project
	CreateProject(projectName string) (Project, error)

	// Creates a new activity in the store and returns the newly created activity.
	// CreateActivity returns ErrProjectNotFound if projectId does not exist.
	CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (Activity, error)

	// Creates a new dependency table row and returns the newly created row.
	// CreateTableRow returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId does not exist.
	CreateTableRow(projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (TableRow, error)

	// Creates a new dependency table row with a sort rank greater than the maximum sort rank by [TableRowSortRankStep] amount and returns the newly created row.
	// CreateTableRowAtBottom returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId does not exist.
	CreateTableRowAtBottom(projectId uuid.UUID, activityId uuid.UUID) (TableRow, error)

	// Creates a new dependency table row below the row associated with belowActivityId and returns the newly created row.
	// CreateTableRowBelow returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if activityId or belowActivityId does not exist.
	CreateTableRowBelow(projectId uuid.UUID, activityId uuid.UUID, belowActivityId uuid.UUID) (TableRow, error)

	// Creates a new dependency between the successor activity and predecessor activity and returns the newly created dependency.
	// CreateDependency returns ErrProjectNotFound if projectId does not exist,
	// and ErrActivityNotFound if successorActivityId or predecessorActivityId does not exist.
	CreateDependency(projectId uuid.UUID, successorActivityId uuid.UUID, predecessorActivityId uuid.UUID) (Dependency, error)

	// Gets the activities associated with the projectId.
	// GetActivities returns ErrProjectNotFound if projectId does not exist.
	GetActivities(projectId uuid.UUID) ([]Activity, error)

	// Gets the table rows associated with the projectId sorted by sort rank.
	// GetTableRowsSortedByRank returns ErrProjectNotFound if projectId does not exist.
	GetTableRowsSortedByRank(projectId uuid.UUID) ([]TableRow, error)

	// Updates the name of the activity associated with activityId.
	// UpdateActivityName returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityName(activityId uuid.UUID, activityName string) (Activity, error)

	// Updates the duration of the activity associated with the activityId.
	// UpdateActivityDuration returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityDuration(activityId uuid.UUID, duration time.Duration) (Activity, error)

	// Updates the dependencies of the activity associated with the activityId.
	// UpdateActivityDependencies returns ErrActivityNotFound if activityId does not exist.
	UpdateActivityDependencies(activityId uuid.UUID, dependencyActivityIds []uuid.UUID) ([]Activity, error)
}
