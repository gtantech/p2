package routes

import (
	"uuid"

	"github.com/gtantech/p2/internal/models/viewmodels"
)

type StoreView interface {
	GetDependencyTableByProjectId(projectId uuid.UUID) (viewmodels.Table, error)
}
