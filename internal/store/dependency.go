package store

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/gtantech/p2/internal/db"
)

type RelationshipType string

var ErrDependencyNotFound = errors.New("dependency not found")

const (
	SS RelationshipType = "SS" // start to start relationship
	FS RelationshipType = "FS" // finish to start relationship
	SF RelationshipType = "SF" // start to finish relationship
	FF RelationshipType = "FF" // finish to finish relationship
)

type Dependency struct {
	ID                    uuid.UUID
	ProjectID             uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

func newDependency(id uuid.UUID, projectId uuid.UUID, relationship RelationshipType, predecessorActivityId uuid.UUID, successorActivityId uuid.UUID) Dependency {
	return Dependency{
		ID:                    id,
		ProjectID:             projectId,
		Relationship:          relationship,
		PredecessorActivityID: predecessorActivityId,
		SuccessorActivityID:   successorActivityId,
	}
}

type CreateDepdencencyParams struct {
	ProjectID             uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

func (c *CreateDepdencencyParams) toDbInsertDependencyParams(dependencyId uuid.UUID) db.InsertDependencyParams {
	return db.InsertDependencyParams{
		ID:                    dependencyId.String(),
		ProjectID:             c.ProjectID.String(),
		Relationship:          string(c.Relationship),
		PredecessorActivityID: c.PredecessorActivityID.String(),
		SuccessorActivityID:   c.SuccessorActivityID.String(),
	}
}

type UpdateDepdencencyParams struct {
	ID                    uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

func (u *UpdateDepdencencyParams) toDbUpdateDependencyParams() db.UpdateDependencyParams {
	return db.UpdateDependencyParams{
		Relationship:          string(u.Relationship),
		PredecessorActivityID: u.PredecessorActivityID.String(),
		SuccessorActivityID:   u.SuccessorActivityID.String(),
		ID:                    u.ID.String(),
	}
}

type DependencyStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (Dependency, error)
	GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]Dependency, error)
	GetByPredecessor(ctx context.Context, predecessorId uuid.UUID) ([]Dependency, error)
	GetBySuccessor(ctx context.Context, successorId uuid.UUID) ([]Dependency, error)
	Create(ctx context.Context, params CreateDepdencencyParams) (Dependency, error)
	Update(ctx context.Context, params UpdateDepdencencyParams) (Dependency, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type dependencyDbStore struct {
	queries *db.Queries
}

// GetByProjectID implements [DependencyStore].
func (d *dependencyDbStore) GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]Dependency, error) {
	data, err := d.queries.FindAllDependenciesByProject(ctx, projectId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return []Dependency{}, ErrDependencyNotFound
		}
		return []Dependency{}, err
	}

	dependencies := make([]Dependency, len(data))
	for i, d := range data {
		dependencies[i] = newDependency(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), RelationshipType(d.Relationship), uuid.MustParse(d.PredecessorActivityID), uuid.MustParse(d.SuccessorActivityID))
	}
	return dependencies, nil
}

// GetByID implements [DependencyStore].
func (d *dependencyDbStore) GetByID(ctx context.Context, id uuid.UUID) (Dependency, error) {
	data, err := d.queries.FindDependencyById(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return Dependency{}, ErrDependencyNotFound
		}
		return Dependency{}, err
	}
	return newDependency(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), RelationshipType(data.Relationship), uuid.MustParse(data.PredecessorActivityID), uuid.MustParse(data.SuccessorActivityID)), nil
}

// GetByPredecessor implements [DependencyStore].
func (d *dependencyDbStore) GetByPredecessor(ctx context.Context, predecessorId uuid.UUID) ([]Dependency, error) {
	data, err := d.queries.FindAllDependenciesByPredecessor(ctx, predecessorId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return []Dependency{}, ErrDependencyNotFound
		}
		return []Dependency{}, err
	}

	dependencies := make([]Dependency, len(data))
	for i, d := range data {
		dependencies[i] = newDependency(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), RelationshipType(d.Relationship), uuid.MustParse(d.PredecessorActivityID), uuid.MustParse(d.SuccessorActivityID))
	}
	return dependencies, nil
}

// GetBySuccessor implements [DependencyStore].
func (d *dependencyDbStore) GetBySuccessor(ctx context.Context, successorId uuid.UUID) ([]Dependency, error) {
	data, err := d.queries.FindAllDependenciesBySuccessor(ctx, successorId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return []Dependency{}, ErrDependencyNotFound
		}
		return []Dependency{}, err
	}

	dependencies := make([]Dependency, len(data))
	for i, d := range data {
		dependencies[i] = newDependency(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), RelationshipType(d.Relationship), uuid.MustParse(d.PredecessorActivityID), uuid.MustParse(d.SuccessorActivityID))
	}
	return dependencies, nil
}

// Create implements [DependencyStore].
func (d *dependencyDbStore) Create(ctx context.Context, params CreateDepdencencyParams) (Dependency, error) {
	data, err := d.queries.InsertDependency(ctx, params.toDbInsertDependencyParams(uuid.NewV7()))
	if err != nil {
		return Dependency{}, err
	}
	return newDependency(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), RelationshipType(data.Relationship), uuid.MustParse(data.PredecessorActivityID), uuid.MustParse(data.SuccessorActivityID)), nil
}

// Delete implements [DependencyStore].
func (d *dependencyDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return d.queries.DeleteDependency(ctx, id.String())
}

// Update implements [DependencyStore].
func (d *dependencyDbStore) Update(ctx context.Context, params UpdateDepdencencyParams) (Dependency, error) {
	data, err := d.queries.UpdateDependency(ctx, params.toDbUpdateDependencyParams())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Dependency doesn't exist
			return Dependency{}, ErrDependencyNotFound
		}
		return Dependency{}, err
	}
	return newDependency(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), RelationshipType(data.Relationship), uuid.MustParse(data.PredecessorActivityID), uuid.MustParse(data.SuccessorActivityID)), nil
}

func NewDependencyStoreFromDb(queries *db.Queries) *dependencyDbStore {
	return &dependencyDbStore{
		queries: queries,
	}
}

var _ DependencyStore = (*dependencyDbStore)(nil) //ensures dependencyDbStore implements DependencyStore at compile time
