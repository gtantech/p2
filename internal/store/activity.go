package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/db"
	"github.com/gtantech/p2/internal/models"
)

var ErrActivityNotFound = errors.New("activity not found")

func newActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) models.Activity {
	return models.Activity{ID: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}

func toDbInsertActivityParams(activityId uuid.UUID, c *models.CreateActivityParams) db.InsertActivityParams {
	return db.InsertActivityParams{
		ID:        activityId.String(),
		ProjectID: c.ProjectID.String(),
		DispName:  c.DisplayName,
		Duration:  int64(c.Duration),
	}
}

func toDbFindAllActivitiesByNameAndProjectParams(g *models.GetActivityByNameAndProjectParams) db.FindAllActivitiesByNameAndProjectParams {
	return db.FindAllActivitiesByNameAndProjectParams{
		DispName:  g.DisplayName,
		ProjectID: g.ProjectID.String(),
	}
}

func toUpdateActivityParams(u *models.UpdateActivityParams) db.UpdateActivityParams {
	return db.UpdateActivityParams{
		DispName: u.DisplayName,
		Duration: int64(u.Duration),
		ID:       u.Id.String(),
	}
}

type activityDbStore struct {
	queries *db.Queries
}

func (a *activityDbStore) Create(ctx context.Context, params models.CreateActivityParams) (models.Activity, error) {
	data, err := a.queries.InsertActivity(ctx, toDbInsertActivityParams(uuid.NewV7(), &params))
	if err != nil {
		return models.Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

func (a *activityDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return a.queries.DeleteActivity(ctx, id.String())
}

func (a *activityDbStore) GetByID(ctx context.Context, id uuid.UUID) (models.Activity, error) {
	data, err := a.queries.FindActivityById(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Activity doesn't exist
			return models.Activity{}, ErrActivityNotFound
		}
		return models.Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

func (a *activityDbStore) GetByNameAndProject(ctx context.Context, params models.GetActivityByNameAndProjectParams) ([]models.Activity, error) {
	data, err := a.queries.FindAllActivitiesByNameAndProject(ctx, toDbFindAllActivitiesByNameAndProjectParams(&params))
	if err != nil {
		return []models.Activity{}, err
	}
	if len(data) == 0 {
		// Activity doesn't exist
		return []models.Activity{}, ErrActivityNotFound
	}
	activities := make([]models.Activity, len(data))
	for i, d := range data {
		activities[i] = newActivity(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), d.DispName, time.Duration(d.Duration))
	}
	return activities, nil
}

func (a *activityDbStore) GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]models.Activity, error) {
	data, err := a.queries.FindAllActivitiesByProject(ctx, projectId.String())
	if err != nil {
		return []models.Activity{}, err
	}
	if len(data) == 0 {
		// Activity doesn't exist
		return []models.Activity{}, ErrActivityNotFound
	}
	activities := make([]models.Activity, len(data))
	for i, d := range data {
		activities[i] = newActivity(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), d.DispName, time.Duration(d.Duration))
	}
	return activities, nil
}

func (a *activityDbStore) Update(ctx context.Context, params models.UpdateActivityParams) (models.Activity, error) {
	data, err := a.queries.UpdateActivity(ctx, toUpdateActivityParams(&params))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Activity doesn't exist
			return models.Activity{}, ErrActivityNotFound
		}
		return models.Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

func NewActivityStoreFromDb(queries *db.Queries) *activityDbStore {
	return &activityDbStore{
		queries: queries,
	}
}
