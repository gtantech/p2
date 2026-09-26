package store

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/gtantech/p2/internal/db"
	"github.com/gtantech/p2/internal/models"
)

var ErrDependencyNotFound = errors.New("dependency not found")

func newDependency(id uuid.UUID, projectId uuid.UUID, relationship models.RelationshipType, predecessorActivityId uuid.UUID, successorActivityId uuid.UUID) models.Dependency {
	return models.Dependency{
		ID:                    id,
		ProjectID:             projectId,
		Relationship:          relationship,
		PredecessorActivityID: predecessorActivityId,
		SuccessorActivityID:   successorActivityId,
	}
}

func toDbInsertDependencyParams(dependencyId uuid.UUID, c *models.CreateDepdencencyParams) db.InsertDependencyParams {
	return db.InsertDependencyParams{
		ID:                    dependencyId.String(),
		ProjectID:             c.ProjectID.String(),
		Relationship:          string(c.Relationship),
		PredecessorActivityID: c.PredecessorActivityID.String(),
		SuccessorActivityID:   c.SuccessorActivityID.String(),
	}
}
func toDbUpdateDependencyParams(u *models.UpdateDepdencencyParams) db.UpdateDependencyParams {
	return db.UpdateDependencyParams{
		Relationship:          string(u.Relationship),
		PredecessorActivityID: u.PredecessorActivityID.String(),
		SuccessorActivityID:   u.SuccessorActivityID.String(),
		ID:                    u.ID.String(),
	}
}

type dependencyDbStore struct {
	queries *db.Queries
}

func (d *dependencyDbStore) GetPredecessorNamesBySuccessor(ctx context.Context, successorId uuid.UUID) ([]models.GetPredecessorNamesBySuccessorResult, error) {
	data, err := d.queries.FindAllPredecessorNamesBySuccessor(ctx, successorId.String())
	if err != nil {
		return []models.GetPredecessorNamesBySuccessorResult{}, err
	}
	if len(data) == 0 {
		// Dependency doesn't exist
		return []models.GetPredecessorNamesBySuccessorResult{}, ErrDependencyNotFound
	}
	predecessors := make([]models.GetPredecessorNamesBySuccessorResult, len(data))
	for i, d := range data {
		predecessors[i] = models.GetPredecessorNamesBySuccessorResult{DependencyID: uuid.MustParse(d.DependencyID), Relationship: models.RelationshipType(d.Relationship), PredecessorActivityID: uuid.MustParse(d.PredecessorActivityID), PredecessorActivityName: d.PredecessorActivityName}
	}
	return predecessors, nil
}

func (d *dependencyDbStore) GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]models.Dependency, error) {
	data, err := d.queries.FindAllDependenciesByProject(ctx, projectId.String())
	if err != nil {
		return []models.Dependency{}, err
	}
	if len(data) == 0 {
		// Dependency doesn't exist
		return []models.Dependency{}, ErrDependencyNotFound
	}
	dependencies := make([]models.Dependency, len(data))
	for i, d := range data {
		dependencies[i] = newDependency(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), models.RelationshipType(d.Relationship), uuid.MustParse(d.PredecessorActivityID), uuid.MustParse(d.SuccessorActivityID))
	}
	return dependencies, nil
}

func (d *dependencyDbStore) GetByID(ctx context.Context, id uuid.UUID) (models.Dependency, error) {
	data, err := d.queries.FindDependencyById(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return models.Dependency{}, ErrDependencyNotFound
		}
		return models.Dependency{}, err
	}
	return newDependency(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), models.RelationshipType(data.Relationship), uuid.MustParse(data.PredecessorActivityID), uuid.MustParse(data.SuccessorActivityID)), nil
}

func (d *dependencyDbStore) GetByPredecessor(ctx context.Context, predecessorId uuid.UUID) ([]models.Dependency, error) {
	data, err := d.queries.FindAllDependenciesByPredecessor(ctx, predecessorId.String())
	if err != nil {
		return []models.Dependency{}, err
	}
	if len(data) == 0 {
		// Dependency doesn't exist
		return []models.Dependency{}, ErrDependencyNotFound
	}
	dependencies := make([]models.Dependency, len(data))
	for i, d := range data {
		dependencies[i] = newDependency(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), models.RelationshipType(d.Relationship), uuid.MustParse(d.PredecessorActivityID), uuid.MustParse(d.SuccessorActivityID))
	}
	return dependencies, nil
}

func (d *dependencyDbStore) GetBySuccessor(ctx context.Context, successorId uuid.UUID) ([]models.Dependency, error) {
	data, err := d.queries.FindAllDependenciesBySuccessor(ctx, successorId.String())
	if err != nil {
		return []models.Dependency{}, err
	}
	if len(data) == 0 {
		// Dependency doesn't exist
		return []models.Dependency{}, ErrDependencyNotFound
	}
	dependencies := make([]models.Dependency, len(data))
	for i, d := range data {
		dependencies[i] = newDependency(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), models.RelationshipType(d.Relationship), uuid.MustParse(d.PredecessorActivityID), uuid.MustParse(d.SuccessorActivityID))
	}
	return dependencies, nil
}

func (d *dependencyDbStore) Create(ctx context.Context, params models.CreateDepdencencyParams) (models.Dependency, error) {
	data, err := d.queries.InsertDependency(ctx, toDbInsertDependencyParams(uuid.NewV7(), &params))
	if err != nil {
		return models.Dependency{}, err
	}
	return newDependency(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), models.RelationshipType(data.Relationship), uuid.MustParse(data.PredecessorActivityID), uuid.MustParse(data.SuccessorActivityID)), nil
}

func (d *dependencyDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return d.queries.DeleteDependency(ctx, id.String())
}

func (d *dependencyDbStore) Update(ctx context.Context, params models.UpdateDepdencencyParams) (models.Dependency, error) {
	data, err := d.queries.UpdateDependency(ctx, toDbUpdateDependencyParams(&params))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return models.Dependency{}, ErrDependencyNotFound
		}
		return models.Dependency{}, err
	}
	return newDependency(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), models.RelationshipType(data.Relationship), uuid.MustParse(data.PredecessorActivityID), uuid.MustParse(data.SuccessorActivityID)), nil
}

func NewDependencyStoreFromDb(queries *db.Queries) *dependencyDbStore {
	return &dependencyDbStore{
		queries: queries,
	}
}
