package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/db"
)

var ErrActivityNotFound = errors.New("activity not found")

type Activity struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

func newActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) Activity {
	return Activity{ID: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}

type CreateActivityParams struct {
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

func toDbInsertActivityParams(activityId uuid.UUID, c *CreateActivityParams) db.InsertActivityParams {
	return db.InsertActivityParams{
		ID:        activityId.String(),
		ProjectID: c.ProjectID.String(),
		DispName:  c.DisplayName,
		Duration:  int64(c.Duration),
	}
}

type GetActivityByNameAndProjectParams struct {
	ProjectID   uuid.UUID
	DisplayName string
}

func toDbFindAllActivitiesByNameAndProjectParams(g *GetActivityByNameAndProjectParams) db.FindAllActivitiesByNameAndProjectParams {
	return db.FindAllActivitiesByNameAndProjectParams{
		DispName:  g.DisplayName,
		ProjectID: g.ProjectID.String(),
	}
}

type UpdateActivityParams struct {
	Id          uuid.UUID
	DisplayName string
	Duration    time.Duration
}

func toUpdateActivityParams(u *UpdateActivityParams) db.UpdateActivityParams {
	return db.UpdateActivityParams{
		DispName: u.DisplayName,
		Duration: int64(u.Duration),
		ID:       u.Id.String(),
	}
}

type activityDbStore struct {
	queries *db.Queries
}

func (a *activityDbStore) Create(ctx context.Context, params CreateActivityParams) (Activity, error) {
	data, err := a.queries.InsertActivity(ctx, toDbInsertActivityParams(uuid.NewV7(), &params))
	if err != nil {
		return Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

func (a *activityDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return a.queries.DeleteActivity(ctx, id.String())
}

func (a *activityDbStore) GetByID(ctx context.Context, id uuid.UUID) (Activity, error) {
	data, err := a.queries.FindActivityById(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Activity doesn't exist
			return Activity{}, ErrActivityNotFound
		}
		return Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

func (a *activityDbStore) GetByNameAndProject(ctx context.Context, params GetActivityByNameAndProjectParams) ([]Activity, error) {
	data, err := a.queries.FindAllActivitiesByNameAndProject(ctx, toDbFindAllActivitiesByNameAndProjectParams(&params))
	if err != nil {
		return []Activity{}, err
	}
	if len(data) == 0 {
		// Activity doesn't exist
		return []Activity{}, ErrActivityNotFound
	}
	activities := make([]Activity, len(data))
	for i, d := range data {
		activities[i] = newActivity(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), d.DispName, time.Duration(d.Duration))
	}
	return activities, nil
}

func (a *activityDbStore) GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]Activity, error) {
	data, err := a.queries.FindAllActivitiesByProject(ctx, projectId.String())
	if err != nil {
		return []Activity{}, err
	}
	if len(data) == 0 {
		// Activity doesn't exist
		return []Activity{}, ErrActivityNotFound
	}
	activities := make([]Activity, len(data))
	for i, d := range data {
		activities[i] = newActivity(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), d.DispName, time.Duration(d.Duration))
	}
	return activities, nil
}

func (a *activityDbStore) Update(ctx context.Context, params UpdateActivityParams) (Activity, error) {
	data, err := a.queries.UpdateActivity(ctx, toUpdateActivityParams(&params))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Activity doesn't exist
			return Activity{}, ErrActivityNotFound
		}
		return Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

func NewActivityStoreFromDb(queries *db.Queries) *activityDbStore {
	return &activityDbStore{
		queries: queries,
	}
}
