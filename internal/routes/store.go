package routes

import (
	"context"
	"uuid"

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
	GetByID(ctx context.Context, id uuid.UUID) (store.Project, error)
	GetByName(ctx context.Context, search string) ([]store.Project, error)
	GetProjects(ctx context.Context) ([]store.Project, error)
	Create(ctx context.Context, params store.CreateProjectParams) (store.Project, error)
	Update(ctx context.Context, params store.UpdateProjectParams) (store.Project, error)
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
