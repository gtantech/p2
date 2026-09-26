package store

import (
	"context"
	"database/sql"
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

type CreateProjectParams struct {
	DisplayName string
}

func (c *CreateProjectParams) toDbInsertProjectParams(projectId uuid.UUID) db.InsertProjectParams {
	return db.InsertProjectParams{ID: projectId.String(), DispName: c.DisplayName}
}

type UpdateProjectParams struct {
	Id          uuid.UUID
	DisplayName string
}

func (u *UpdateProjectParams) toDbUpdateProjectParams() db.UpdateProjectParams {
	return db.UpdateProjectParams{ID: u.Id.String(), DispName: u.DisplayName}
}

type projectDbStore struct {
	queries *db.Queries
}

// Update implements [ProjectStore].
func (p *projectDbStore) Update(ctx context.Context, params UpdateProjectParams) (Project, error) {
	data, err := p.queries.UpdateProject(ctx, params.toDbUpdateProjectParams())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Project doesn't exist
			return Project{}, ErrProjectNotFound
		}
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
	if len(data) == 0 {
		// Project doesn't exist
		return []Project{}, ErrProjectNotFound
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
	if len(data) == 0 {
		// Project doesn't exist
		return []Project{}, ErrProjectNotFound
	}
	projects := make([]Project, len(data))
	for i, d := range data {
		projects[i] = newProject(uuid.MustParse(d.ID), d.DispName)
	}
	return projects, nil
}

// Create implements [ProjectStore].
func (p *projectDbStore) Create(ctx context.Context, params CreateProjectParams) (Project, error) {
	data, err := p.queries.InsertProject(ctx, params.toDbInsertProjectParams(uuid.NewV7()))
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
		if errors.Is(err, sql.ErrNoRows) {
			// Project doesn't exist
			return Project{}, ErrProjectNotFound
		}
		return Project{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

func NewProjectStoreFromDb(queries *db.Queries) *projectDbStore {
	return &projectDbStore{
		queries: queries,
	}
}
