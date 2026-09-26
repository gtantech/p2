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
	Create(ctx context.Context, params models.StoreCreateProjectParams) (models.StoreProject, error)
	Update(ctx context.Context, params models.StoreUpdateProjectParams) (models.StoreProject, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DependencyStoreService interface {
	GetByID(ctx context.Context, id uuid.UUID) (models.StoreDependency, error)
	GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]models.StoreDependency, error)
	GetByPredecessor(ctx context.Context, predecessorId uuid.UUID) ([]models.StoreDependency, error)
	GetBySuccessor(ctx context.Context, successorId uuid.UUID) ([]models.StoreDependency, error)
	GetPredecessorNamesBySuccessor(ctx context.Context, successorId uuid.UUID) ([]models.StoreGetPredecessorNamesBySuccessorResult, error)
	Create(ctx context.Context, params models.StoreCreateDepdencencyParams) (models.StoreDependency, error)
	Update(ctx context.Context, params models.StoreUpdateDepdencencyParams) (models.StoreDependency, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
