package routes

import (
	"context"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models/storemodels"
	"github.com/gtantech/p2/internal/models/viewmodels"
)

type StoreView interface {
	GetDependencyTableByProjectId(projectId uuid.UUID, ctx context.Context) (viewmodels.Table, error)
	CreateEmptyDependencyTableRow(projectId uuid.UUID, sortRank int64, ctx context.Context) (viewmodels.TableRow, error)
	CreateDependencyTableRow(projectId uuid.UUID, activityName string, sortRank int64, ctx context.Context) (viewmodels.TableRow, error)
}

type Store interface {
	UpdateActivityName(activityId uuid.UUID, activityName string, ctx context.Context) error
	UpdateActivityDuration(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error
	GetActivitiesByProjectId(projectId uuid.UUID, ctx context.Context) ([]storemodels.Activity, error)
	GetLastDependencyTableRowByProjectId(projectId uuid.UUID, ctx context.Context) (storemodels.TableRow, error)
	GetPredecessorActivityNamesBySuccessorId(successorId uuid.UUID, ctx context.Context) ([]storemodels.ActivityNameWithId, error)
	DeleteDependencyByProjectPredecessorSuccessorId(projectId uuid.UUID, predecessorId uuid.UUID, successorId uuid.UUID, ctx context.Context) error
}
