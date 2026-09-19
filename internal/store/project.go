package store

import (
	"context"
	"errors"
	"uuid"

	"github.com/gtantech/p2/internal/db"
)

var ErrProjectNotFound = errors.New("project not found")

type Project struct {
	ID          uuid.UUID
	DisplayName string
}

func newProject(id uuid.UUID, displayName string) Project {
	return Project{ID: id, DisplayName: displayName}
}

type ProjectStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (Project, error)
	GetByName(ctx context.Context, search string) ([]Project, error)
	GetProjects(ctx context.Context) ([]Project, error)
	Create(ctx context.Context, params CreateProjectDbParams) (Project, error)
	Update(ctx context.Context, params UpdateProjectDbParams) (Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateProjectDbParams struct {
	DisplayName string
}

type UpdateProjectDbParams struct {
	Id          uuid.UUID
	DisplayName string
}

type projectDbStore struct {
	queries *db.Queries
}

// Update implements [ProjectStore].
func (p *projectDbStore) Update(ctx context.Context, params UpdateProjectDbParams) (Project, error) {
	data, err := p.queries.UpdateProject(ctx, db.UpdateProjectParams{ID: params.Id.String(), DispName: params.DisplayName})
	if err != nil {
		return Project{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

// GetProjects implements [ProjectStore].
func (p *projectDbStore) GetProjects(ctx context.Context) ([]Project, error) {
	data, err := p.queries.FindAllProjects(ctx)
	if err != nil {
		return []Project{}, err
	}
	projects := make([]Project, len(data))
	for i, d := range data {
		projects[i] = newProject(uuid.MustParse(d.ID), d.DispName)
	}
	return projects, nil
}

// GetByName implements [ProjectStore].
func (p *projectDbStore) GetByName(ctx context.Context, search string) ([]Project, error) {
	data, err := p.queries.FindProjectByName(ctx, search)
	if err != nil {
		return []Project{}, err
	}
	projects := make([]Project, len(data))
	for i, d := range data {
		projects[i] = newProject(uuid.MustParse(d.ID), d.DispName)
	}
	return projects, nil
}

// Create implements [ProjectStore].
func (p *projectDbStore) Create(ctx context.Context, params CreateProjectDbParams) (Project, error) {
	data, err := p.queries.InsertProject(ctx, db.InsertProjectParams{ID: uuid.NewV7().String(), DispName: params.DisplayName})
	if err != nil {
		return Project{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

// Delete implements [ProjectStore].
func (p *projectDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return p.queries.DeleteProject(ctx, id.String())
}

// GetByID implements [ProjectStore].
func (p *projectDbStore) GetByID(ctx context.Context, id uuid.UUID) (Project, error) {
	data, err := p.queries.FindProjectById(ctx, id.String())
	if err != nil {
		return Project{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

func NewProjectDbStore(queries *db.Queries) *projectDbStore {
	return &projectDbStore{
		queries: queries,
	}
}

var _ ProjectStore = (*projectDbStore)(nil) //ensures projectStore implements ProjectStore at compile time
