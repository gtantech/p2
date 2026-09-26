package routes

import (
	"context"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/store"
)

type StoreService interface {
	Activity() ActivityStoreService
	Project() ProjectStoreService
	Dependency() DependencyStoreService
}

type ActivityStoreService interface {
	GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]store.Activity, error)
	GetByNameAndProject(ctx context.Context, params store.GetActivityByNameAndProjectParams) ([]store.Activity, error)
	GetByID(ctx context.Context, id uuid.UUID) (store.Activity, error)
	Create(ctx context.Context, params store.CreateActivityParams) (store.Activity, error)
	Update(ctx context.Context, params store.UpdateActivityParams) (store.Activity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectStoreService interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.StoreProject, error)
	GetByName(ctx context.Context, search string) ([]models.StoreProject, error)
	GetProjects(ctx context.Context) ([]models.StoreProject, error)
	Create(ctx context.Context, params models.CreateProjectParams) (models.StoreProject, error)
	Update(ctx context.Context, params models.UpdateProjectParams) (models.StoreProject, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DependencyStoreService interface {
	GetByID(ctx context.Context, id uuid.UUID) (store.Dependency, error)
	GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]store.Dependency, error)
	GetByPredecessor(ctx context.Context, predecessorId uuid.UUID) ([]store.Dependency, error)
	GetBySuccessor(ctx context.Context, successorId uuid.UUID) ([]store.Dependency, error)
	GetPredecessorNamesBySuccessor(ctx context.Context, successorId uuid.UUID) ([]store.GetPredecessorNamesBySuccessorResult, error)
	Create(ctx context.Context, params store.CreateDepdencencyParams) (store.Dependency, error)
	Update(ctx context.Context, params store.UpdateDepdencencyParams) (store.Dependency, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
