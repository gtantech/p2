package store

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/gtantech/p2/internal/db"
	"github.com/gtantech/p2/internal/models"
)

var ErrProjectNotFound = errors.New("project not found")

func newProject(id uuid.UUID, displayName string) models.StoreProject {
	return models.StoreProject{ID: id, DisplayName: displayName}
}

func toDbInsertProjectParams(projectId uuid.UUID, c *models.StoreCreateProjectParams) db.InsertProjectParams {
	return db.InsertProjectParams{ID: projectId.String(), DispName: c.DisplayName}
}

func toDbUpdateProjectParams(u *models.StoreUpdateProjectParams) db.UpdateProjectParams {
	return db.UpdateProjectParams{ID: u.Id.String(), DispName: u.DisplayName}
}

type projectDbStore struct {
	queries *db.Queries
}

func (p *projectDbStore) Update(ctx context.Context, params models.StoreUpdateProjectParams) (models.StoreProject, error) {
	data, err := p.queries.UpdateProject(ctx, toDbUpdateProjectParams(&params))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Project doesn't exist
			return models.StoreProject{}, ErrProjectNotFound
		}
		return models.StoreProject{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

func (p *projectDbStore) GetProjects(ctx context.Context) ([]models.StoreProject, error) {
	data, err := p.queries.FindAllProjects(ctx)
	if err != nil {
		return []models.StoreProject{}, err
	}
	if len(data) == 0 {
		// Project doesn't exist
		return []models.StoreProject{}, ErrProjectNotFound
	}
	projects := make([]models.StoreProject, len(data))
	for i, d := range data {
		projects[i] = newProject(uuid.MustParse(d.ID), d.DispName)
	}
	return projects, nil
}

func (p *projectDbStore) GetByName(ctx context.Context, search string) ([]models.StoreProject, error) {
	data, err := p.queries.FindProjectByName(ctx, search)
	if err != nil {
		return []models.StoreProject{}, err
	}
	if len(data) == 0 {
		// Project doesn't exist
		return []models.StoreProject{}, ErrProjectNotFound
	}
	projects := make([]models.StoreProject, len(data))
	for i, d := range data {
		projects[i] = newProject(uuid.MustParse(d.ID), d.DispName)
	}
	return projects, nil
}

func (p *projectDbStore) Create(ctx context.Context, params models.StoreCreateProjectParams) (models.StoreProject, error) {
	data, err := p.queries.InsertProject(ctx, toDbInsertProjectParams(uuid.NewV7(), &params))
	if err != nil {
		return models.StoreProject{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

func (p *projectDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return p.queries.DeleteProject(ctx, id.String())
}

func (p *projectDbStore) GetByID(ctx context.Context, id uuid.UUID) (models.StoreProject, error) {
	data, err := p.queries.FindProjectById(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Project doesn't exist
			return models.StoreProject{}, ErrProjectNotFound
		}
		return models.StoreProject{}, err
	}
	return newProject(uuid.MustParse(data.ID), data.DispName), nil
}

func NewProjectStoreFromDb(queries *db.Queries) *projectDbStore {
	return &projectDbStore{
		queries: queries,
	}
}
