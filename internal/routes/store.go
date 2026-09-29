package routes

import (
	"context"
	"uuid"

	"github.com/gtantech/p2/internal/models/viewmodels"
)

type StoreView interface {
	GetDependencyTableByProjectId(projectId uuid.UUID, ctx context.Context) (viewmodels.Table, error)
}
