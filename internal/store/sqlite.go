package store

import (
	"context"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/database/sqlitedb"
	"github.com/gtantech/p2/internal/models"
)

var _ models.Model = (*SqliteStore)(nil) //ensures SqliteStore implements models.Model at compile time

func NewSqliteStore(queries sqlitedb.Queries) *SqliteStore {
	return &SqliteStore{queries: queries}
}

type SqliteStore struct {
	queries sqlitedb.Queries
}

// CreateActivity implements [models.Model].
func (s *SqliteStore) CreateActivity(ctx context.Context, projectId uuid.UUID, activityName string, duration time.Duration) (models.Activity, error) {
	panic("unimplemented")
}

// CreateDependency implements [models.Model].
func (s *SqliteStore) CreateDependency(ctx context.Context, projectId uuid.UUID, successorActivityId uuid.UUID, predecessorActivityId uuid.UUID) (models.Dependency, error) {
	panic("unimplemented")
}

// CreateProject implements [models.Model].
func (s *SqliteStore) CreateProject(ctx context.Context, projectName string) (models.Project, error) {
	panic("unimplemented")
}

// CreateTableRow implements [models.Model].
func (s *SqliteStore) CreateTableRow(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (models.TableRow, error) {
	panic("unimplemented")
}

// CreateTableRowAtBottom implements [models.Model].
func (s *SqliteStore) CreateTableRowAtBottom(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID) (models.TableRow, error) {
	panic("unimplemented")
}

// CreateTableRowBelow implements [models.Model].
func (s *SqliteStore) CreateTableRowBelow(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, belowActivityId uuid.UUID) (models.TableRow, error) {
	panic("unimplemented")
}

// GetActivities implements [models.Model].
func (s *SqliteStore) GetActivities(ctx context.Context, projectId uuid.UUID) ([]models.Activity, error) {
	panic("unimplemented")
}

// GetTableRowsSortedByRank implements [models.Model].
func (s *SqliteStore) GetTableRowsSortedByRank(ctx context.Context, projectId uuid.UUID) ([]models.TableRow, error) {
	panic("unimplemented")
}

// UpdateActivityDependencies implements [models.Model].
func (s *SqliteStore) UpdateActivityDependencies(ctx context.Context, activityId uuid.UUID, dependencyActivityIds []uuid.UUID) ([]models.Activity, error) {
	panic("unimplemented")
}

// UpdateActivityDuration implements [models.Model].
func (s *SqliteStore) UpdateActivityDuration(ctx context.Context, activityId uuid.UUID, duration time.Duration) (models.Activity, error) {
	panic("unimplemented")
}

// UpdateActivityName implements [models.Model].
func (s *SqliteStore) UpdateActivityName(ctx context.Context, activityId uuid.UUID, activityName string) (models.Activity, error) {
	panic("unimplemented")
}
