package store

import (
	"context"
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

type CreateActivityDbParams struct {
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type GetByNameAndProjectParams struct {
	ProjectID   uuid.UUID
	DisplayName string
}

type UpdateActivityDbParams struct {
	Id          uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type ActivityStore interface {
	GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]Activity, error)
	GetByNameAndProject(ctx context.Context, params GetByNameAndProjectParams) ([]Activity, error)
	GetByID(ctx context.Context, id uuid.UUID) (Activity, error)
	Create(ctx context.Context, params CreateActivityDbParams) (Activity, error)
	Update(ctx context.Context, params UpdateActivityDbParams) (Activity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type activityDbStore struct {
	queries *db.Queries
}

// Create implements [ActivityStore].
func (a *activityDbStore) Create(ctx context.Context, params CreateActivityDbParams) (Activity, error) {
	data, err := a.queries.InsertActivity(ctx, db.InsertActivityParams{
		ID:        uuid.NewV7().String(),
		ProjectID: params.ProjectID.String(),
		DispName:  params.DisplayName,
		Duration:  int64(params.Duration),
	})
	if err != nil {
		return Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

// Delete implements [ActivityStore].
func (a *activityDbStore) Delete(ctx context.Context, id uuid.UUID) error {
	return a.queries.DeleteActivity(ctx, id.String())
}

// GetByID implements [ActivityStore].
func (a *activityDbStore) GetByID(ctx context.Context, id uuid.UUID) (Activity, error) {
	data, err := a.queries.FindActivityById(ctx, id.String())
	if err != nil {
		return Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

// GetByNameAndProject implements [ActivityStore].
func (a *activityDbStore) GetByNameAndProject(ctx context.Context, params GetByNameAndProjectParams) ([]Activity, error) {
	data, err := a.queries.FindAllActivitiesByNameAndProject(ctx, db.FindAllActivitiesByNameAndProjectParams{
		DispName:  params.DisplayName,
		ProjectID: params.ProjectID.String(),
	})
	if err != nil {
		return []Activity{}, err
	}
	activities := make([]Activity, len(data))
	for i, d := range data {
		activities[i] = newActivity(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), d.DispName, time.Duration(d.Duration))
	}
	return activities, nil
}

// GetByProjectID implements [ActivityStore].
func (a *activityDbStore) GetByProjectID(ctx context.Context, projectId uuid.UUID) ([]Activity, error) {
	data, err := a.queries.FindAllActivitiesByProject(ctx, projectId.String())
	if err != nil {
		return []Activity{}, err
	}
	activities := make([]Activity, len(data))
	for i, d := range data {
		activities[i] = newActivity(uuid.MustParse(d.ID), uuid.MustParse(d.ProjectID), d.DispName, time.Duration(d.Duration))
	}
	return activities, nil
}

// Update implements [ActivityStore].
func (a *activityDbStore) Update(ctx context.Context, params UpdateActivityDbParams) (Activity, error) {
	data, err := a.queries.UpdateActivity(ctx, db.UpdateActivityParams{
		DispName: params.DisplayName,
		Duration: int64(params.Duration),
		ID:       params.Id.String(),
	})
	if err != nil {
		return Activity{}, err
	}
	return newActivity(uuid.MustParse(data.ID), uuid.MustParse(data.ProjectID), data.DispName, time.Duration(data.Duration)), nil
}

var _ ActivityStore = (*activityDbStore)(nil) //ensures activityDbStore implements ActivityStore at compile time
