package store

import (
	"context"
	"database/sql"
	"errors"
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
	resp, err := s.queries.InsertActivity(ctx, sqlitedb.InsertActivityParams{ID: uuid.NewV7().String(), ProjectID: projectId.String(), DispName: activityName, Duration: int64(duration)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Activity{}, ErrProjectNotFound
		}
		return models.Activity{}, err
	}
	return models.Activity{ActivityId: uuid.MustParse(resp.ID), DisplayName: resp.DispName, Duration: time.Duration(resp.Duration)}, nil
}

// CreateDependency implements [models.Model].
func (s *SqliteStore) CreateDependency(ctx context.Context, projectId uuid.UUID, successorActivityId uuid.UUID, predecessorActivityId uuid.UUID) (models.Dependency, error) {
	resp, err := s.queries.InsertDependency(ctx, sqlitedb.InsertDependencyParams{
		ID:                    uuid.NewV7().String(),
		ProjectID:             projectId.String(),
		Relationship:          string(models.FS),
		PredecessorActivityID: predecessorActivityId.String(),
		SuccessorActivityID:   successorActivityId.String(),
	})
	if err != nil {
		if _, err := s.queries.FindProjectById(ctx, projectId.String()); err != nil {
			return models.Dependency{}, ErrProjectNotFound
		}
		if _, err := s.queries.FindActivityById(ctx, successorActivityId.String()); err != nil {
			return models.Dependency{}, ErrActivityNotFound
		}
		if _, err := s.queries.FindActivityById(ctx, predecessorActivityId.String()); err != nil {
			return models.Dependency{}, ErrActivityNotFound
		}
		return models.Dependency{}, err
	}

	return models.Dependency{
		DependencyId:          uuid.MustParse(resp.ID),
		SuccessorActivityId:   uuid.MustParse(resp.SuccessorActivityID),
		PredecessorActivityId: uuid.MustParse(resp.PredecessorActivityID),
	}, nil
}

// CreateProject implements [models.Model].
func (s *SqliteStore) CreateProject(ctx context.Context, projectName string) (models.Project, error) {
	resp, err := s.queries.InsertProject(ctx, sqlitedb.InsertProjectParams{
		ID:       uuid.NewV7().String(),
		DispName: projectName,
	})
	if err != nil {
		return models.Project{}, err
	}
	return models.Project{
		ProjectId:          uuid.MustParse(resp.ID),
		ProjectDisplayName: resp.DispName,
	}, nil
}

// CreateTableRow implements [models.Model].
func (s *SqliteStore) CreateTableRow(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, rowSortRank int64) (models.TableRow, error) {
	activityOrderingResp, err := s.queries.InsertActivityOrdering(ctx, sqlitedb.InsertActivityOrderingParams{
		ID:                  uuid.NewV7().String(),
		ProjectID:           projectId.String(),
		SuccessorActivityID: activityId.String(),
		SortRank:            rowSortRank,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.TableRow{}, err
		}
		if _, err := s.queries.FindProjectById(ctx, projectId.String()); err != nil && errors.Is(err, sql.ErrNoRows) {
			return models.TableRow{}, ErrProjectNotFound
		}
		if _, err := s.queries.FindActivityById(ctx, activityId.String()); err != nil && errors.Is(err, sql.ErrNoRows) {
			return models.TableRow{}, ErrActivityNotFound
		}
	}
	actvityResp, err := s.queries.FindActivityById(ctx, activityOrderingResp.SuccessorActivityID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.TableRow{}, err
		}
		return models.TableRow{}, ErrActivityNotFound
	}
	dependenciesResp, err := s.queries.FindAllPredecessorNamesBySuccessor(ctx, activityOrderingResp.SuccessorActivityID)
	return models.TableRow{
		Activity: &models.Activity{
			ActivityId:  uuid.MustParse(actvityResp.ID),
			DisplayName: actvityResp.DispName,
			Duration:    time.Duration(actvityResp.Duration),
		},
		ProjectId: uuid.MustParse(actvityResp.ProjectID),
		Dependencies: func() []string {
			dependencies := make([]string, len(dependenciesResp))
			for i, d := range dependenciesResp {
				dependencies[i] = d.PredecessorActivityName
			}
			return dependencies
		}(),
		SortRank: activityOrderingResp.SortRank,
	}, nil
}

// CreateTableRowAtBottom implements [models.Model].
func (s *SqliteStore) CreateTableRowAtBottom(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID) (models.TableRow, error) {
	lastActivityResp, err := s.queries.FindLastActivityByProjectId(ctx, projectId.String())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.TableRow{}, err
		}
		return models.TableRow{}, ErrProjectNotFound
	}
	return s.CreateTableRow(ctx, projectId, activityId, lastActivityResp.SortRank+models.TableRowSortRankStep)
}

// CreateTableRowBelow implements [models.Model].
func (s *SqliteStore) CreateTableRowBelow(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, belowActivityId uuid.UUID) (models.TableRow, error) {
	insertRowRank := int64(0)
	rowCurrentAndBelowResp, err := s.queries.FindCurrentAndNextActivityByProjectIdAndActivityId(ctx, sqlitedb.FindCurrentAndNextActivityByProjectIdAndActivityIdParams{
		ProjectID:           projectId.String(),
		SuccessorActivityID: belowActivityId.String(),
	})
	if err != nil {
		return models.TableRow{}, err
	}
	if len(rowCurrentAndBelowResp) == 0 {
		if _, err := s.queries.FindProjectById(ctx, projectId.String()); err != nil {
			return models.TableRow{}, ErrProjectNotFound
		}
		if _, err := s.queries.FindActivityById(ctx, belowActivityId.String()); err != nil {
			return models.TableRow{}, ErrActivityNotFound
		}
	}
	if len(rowCurrentAndBelowResp) == 1 {
		// end of table
		insertRowRank = rowCurrentAndBelowResp[0].SortRank + models.TableRowSortRankStep
	} else {
		insertRowRank = (rowCurrentAndBelowResp[0].SortRank / 2) + (rowCurrentAndBelowResp[1].SortRank / 2)
	}

	return s.CreateTableRow(ctx, projectId, activityId, insertRowRank)
}

// GetActivities implements [models.Model].
func (s *SqliteStore) GetActivities(ctx context.Context, projectId uuid.UUID) ([]models.Activity, error) {
	resp, err := s.queries.FindAllActivitiesByProject(ctx, projectId.String())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return []models.Activity{}, err
		}
		return []models.Activity{}, ErrProjectNotFound
	}

	return func() []models.Activity {
		activities := make([]models.Activity, len(resp))
		for i, a := range resp {
			activities[i] = models.Activity{
				ActivityId:  uuid.MustParse(a.ID),
				DisplayName: a.DispName,
				Duration:    time.Duration(a.Duration),
			}
		}
		return activities
	}(), nil
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
