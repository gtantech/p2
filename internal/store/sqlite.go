package store

import (
	"context"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/database/sqlite"
	"github.com/gtantech/p2/internal/models"
)

var _ models.Store = (*Sqlite)(nil) //ensures SqliteStore implements models.Store at compile time

func NewSqlite(queries *sqlite.Queries) *Sqlite {
	return &Sqlite{queries: queries}
}

type Sqlite struct {
	queries *sqlite.Queries
}

// CreateProject implements [models.Store].
func (s *Sqlite) CreateProject(ctx context.Context, projectName string) (models.Project, error) {
	resp, err := s.queries.InsertProject(ctx, sqlite.InsertProjectParams{
		ID:       uuid.NewV7().String(),
		DispName: projectName,
	})
	if err != nil {
		return models.Project{}, err
	}
	return models.Project{ProjecId: uuid.MustParse(resp.ID), ProjectName: resp.DispName}, nil
}

// GetHomeProject implements [models.Store].
func (s *Sqlite) GetHomeProject(ctx context.Context) (models.Project, error) {
	homeProjectId := uuid.Max() //TODO remove hard coded value
	resp, err := s.queries.FindProjectById(ctx, homeProjectId.String())
	if err != nil {
		return models.Project{}, ErrProjectNotFound
	}
	return models.Project{
		ProjecId:    uuid.MustParse(resp.ID),
		ProjectName: resp.DispName,
	}, nil
}

// GetTable implements [models.Store].
func (s *Sqlite) GetTable(ctx context.Context, projectId uuid.UUID) (models.Table, error) {
	resp, err := s.queries.FindAllActivitiesByProjectSorted(ctx, projectId.String())
	if err != nil {
		return models.Table{}, err
	}
	if len(resp) == 0 {
		if _, err := s.queries.FindProjectById(ctx, projectId.String()); err != nil {
			return models.Table{}, ErrProjectNotFound
		}
	}

	storeDependencies, err := s.queries.FindAllPredecessorNamesByProjectId(ctx, projectId.String())
	if err != nil {
		return models.Table{}, err
	}

	activityDependencies := make(map[uuid.UUID][]string)
	for _, dependency := range storeDependencies {
		activityId := uuid.MustParse(dependency.SuccessorActivityID)
		dependencyName := dependency.PredecessorActivityName
		activityDependencies[activityId] = append(activityDependencies[activityId], dependencyName)
	}

	tbl := models.Table{
		Rows: func() []models.TableRow {
			tableRows := make([]models.TableRow, len(resp))
			for i, row := range resp {
				tableRows[i] = models.TableRow{
					RowId:           uuid.MustParse(row.RowID),
					ActivityName:    row.DispName,
					DependencyNames: activityDependencies[uuid.MustParse(row.ID)],
					Duration:        time.Duration(row.Duration),
				}
			}
			return tableRows
		}(),
	}
	return tbl, nil
}
