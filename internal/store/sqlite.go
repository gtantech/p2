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

func NewSqliteStore(queries *sqlitedb.Queries) *SqliteStore {
	return &SqliteStore{queries: queries}
}

type SqliteStore struct {
	queries *sqlitedb.Queries
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
	_, err := s.queries.FindProjectById(ctx, projectId.String())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return []models.TableRow{}, err
		}
		return []models.TableRow{}, ErrProjectNotFound
	}
	dependenciesResp, err := s.queries.FindAllPredecessorNamesByProjectId(ctx, projectId.String())
	if err != nil {
		return []models.TableRow{}, err
	}

	activitiesResp, err := s.queries.FindAllActivitiesByProjectSorted(ctx, projectId.String())
	if err != nil {
		return []models.TableRow{}, err
	}

	successorDependencyNames := make(map[uuid.UUID][]string)
	for _, d := range dependenciesResp {
		successorDependencyNames[uuid.MustParse(d.SuccessorActivityID)] = append(successorDependencyNames[uuid.MustParse(d.SuccessorActivityID)], d.PredecessorActivityName)
	}

	tr := make([]models.TableRow, len(activitiesResp))

	for i, a := range activitiesResp {
		activityUuid := uuid.MustParse(a.ID)
		tr[i] = models.TableRow{
			Activity: &models.Activity{
				ActivityId:  activityUuid,
				DisplayName: a.DispName,
				Duration:    time.Duration(a.Duration),
			},
			ProjectId:    projectId,
			Dependencies: successorDependencyNames[activityUuid],
			SortRank:     a.SortRank,
		}
	}
	return tr, nil
}

// UpdateActivityDependencies implements [models.Model].
func (s *SqliteStore) UpdateActivityDependencies(ctx context.Context, projectId uuid.UUID, activityId uuid.UUID, dependencyActivityIds []uuid.UUID) ([]models.Activity, error) {
	dependenciesResp, err := s.queries.FindAllDependenciesBySuccessorId(ctx, activityId.String())
	if err != nil {
		return []models.Activity{}, err
	}

	if len(dependenciesResp) == len(dependencyActivityIds) {
		for i, dr := range dependenciesResp {
			//update only if ids are different
			if dependencyActivityIds[i] != uuid.MustParse(dr.PredecessorActivityID) {
				_, err := s.queries.UpdatePredecessorById(ctx, sqlitedb.UpdatePredecessorByIdParams{
					ID:                    dr.ID,
					PredecessorActivityID: dependencyActivityIds[i].String(),
				})
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return []models.Activity{}, err
					}
					return []models.Activity{}, ErrActivityNotFound
				}
			}
		}
	} else {
		minSize := min(len(dependencyActivityIds), len(dependenciesResp))
		maxSize := max(len(dependencyActivityIds), len(dependenciesResp))
		for i := range minSize {
			dr := dependenciesResp[i]
			if dependencyActivityIds[i] != uuid.MustParse(dr.PredecessorActivityID) {
				_, err := s.queries.UpdatePredecessorById(ctx, sqlitedb.UpdatePredecessorByIdParams{
					ID:                    dr.ID,
					PredecessorActivityID: dependencyActivityIds[i].String(),
				})
				if err != nil {
					if !errors.Is(err, sql.ErrNoRows) {
						return []models.Activity{}, err
					}
					return []models.Activity{}, ErrActivityNotFound
				}
			}
		}

		if maxSize == len(dependenciesResp) {
			//user deleted values
			for i := minSize; i < maxSize; i++ {
				dependencyToDelete := dependenciesResp[i]
				if err := s.queries.DeleteDependencyById(ctx, dependencyToDelete.ID); err != nil {
					return []models.Activity{}, err
				}
			}
		} else {
			//user added values
			for i := minSize; i < maxSize; i++ {
				dependencyIdToAdd := dependencyActivityIds[i]
				_, err := s.queries.InsertDependency(ctx, sqlitedb.InsertDependencyParams{
					ID:                    uuid.NewV7().String(),
					ProjectID:             projectId.String(),
					Relationship:          string(models.FS),
					PredecessorActivityID: dependencyIdToAdd.String(),
					SuccessorActivityID:   activityId.String(),
				})

				if err != nil {
					return []models.Activity{}, err
				}
			}
		}
	}

	dependencyActivitiesResp, err := s.queries.FindAllPredecessorActivitiesBySuccessorActivityId(ctx, activityId.String())

	if err != nil {
		return []models.Activity{}, err
	}

	updatedDependencies := make([]models.Activity, len(dependencyActivitiesResp))

	for i, da := range dependencyActivitiesResp {
		updatedDependencies[i] = models.Activity{
			ActivityId:  uuid.MustParse(da.ID),
			DisplayName: da.DispName,
			Duration:    time.Duration(da.Duration),
		}
	}

	return updatedDependencies, nil

}

// UpdateActivityDuration implements [models.Model].
func (s *SqliteStore) UpdateActivityDuration(ctx context.Context, activityId uuid.UUID, duration time.Duration) (models.Activity, error) {
	a, err := s.queries.UpdateActivityDuration(ctx, sqlitedb.UpdateActivityDurationParams{
		ID:       activityId.String(),
		Duration: int64(duration),
	})

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.Activity{}, err
		}
		return models.Activity{}, ErrActivityNotFound
	}

	return models.Activity{
		ActivityId:  uuid.MustParse(a.ID),
		DisplayName: a.DispName,
		Duration:    time.Duration(a.Duration),
	}, nil
}

// UpdateActivityName implements [models.Model].
func (s *SqliteStore) UpdateActivityName(ctx context.Context, activityId uuid.UUID, activityName string) (models.Activity, error) {
	a, err := s.queries.UpdateActivityName(ctx, sqlitedb.UpdateActivityNameParams{
		ID:       activityId.String(),
		DispName: activityName,
	})

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.Activity{}, err
		}
		return models.Activity{}, ErrActivityNotFound
	}

	return models.Activity{
		ActivityId:  uuid.MustParse(a.ID),
		DisplayName: a.DispName,
		Duration:    time.Duration(a.Duration),
	}, nil
}
