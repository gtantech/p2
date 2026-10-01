package routes

import (
	"context"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models/viewmodels"
)

type StoreView interface {
	GetDependencyTableByProjectId(projectId uuid.UUID, ctx context.Context) (viewmodels.Table, error)
	CreateEmptyDependencyTableRow(projectId uuid.UUID, sortRank int64, ctx context.Context) (viewmodels.TableRow, error)
}

type Store interface {
	UpdateActivityName(activityId uuid.UUID, activityName string, ctx context.Context) error
	UpdateActivityDuration(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error
}
