package sqlite

import (
	"context"
	"database/sql"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/routes"
)

var _ routes.Store = (*Store)(nil) //ensures SqliteStore implements routes.Store at compile time

func NewStore(dbq *DBQueries) *Store {
	return &Store{queries: dbq.queries, db: dbq.db}
}

type Store struct {
	queries *Queries
	db      *sql.DB
}

// CreateProject implements [routes.Store].
func (s *Store) CreateProject(ctx context.Context, projectName string) (models.Project, error) {
	resp, err := s.queries.InsertProject(ctx, InsertProjectParams{
		ID:       uuid.NewV7().String(),
		DispName: projectName,
	})
	if err != nil {
		return models.Project{}, err
	}
	return models.Project{ProjecId: uuid.MustParse(resp.ID), ProjectName: resp.DispName}, nil
}

// GetHomeProject implements [routes.Store].
func (s *Store) GetHomeProject(ctx context.Context) (models.Project, error) {
	homeProjectId := uuid.Max() //TODO remove hard coded value
	resp, err := s.queries.FindProjectById(ctx, homeProjectId.String())
	if err != nil {
		return models.Project{}, routes.ErrProjectNotFound
	}
	return models.Project{
		ProjecId:    uuid.MustParse(resp.ID),
		ProjectName: resp.DispName,
	}, nil
}

// GetTable implements [routes.Store].
func (s *Store) GetTable(ctx context.Context, projectId uuid.UUID) (models.Table, error) {
	resp, err := s.queries.FindAllActivitiesByProjectSorted(ctx, projectId.String())
	if err != nil {
		return models.Table{}, err
	}
	if len(resp) == 0 {
		if _, err := s.queries.FindProjectById(ctx, projectId.String()); err != nil {
			return models.Table{}, routes.ErrProjectNotFound
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
					DependencyNames: activityDependencies[uuid.MustParse(row.ActivityID)],
					Duration:        time.Duration(row.Duration),
				}
			}
			return tableRows
		}(),
	}
	return tbl, nil
}
