package view

import (
	"context"
	"errors"
	"time"
	"uuid"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/store"
)

type View struct {
	store *store.Store
}

var ErrActivityNotFound = errors.New("activity not found")

func NewView(store *store.Store) *View {
	return &View{
		store: store,
	}
}

func (v *View) Home() templ.Component {
	t := NewTable()
	taskA := &Activity{ID: uuid.NewV7(), ProjectID: uuid.NewV7(), DisplayName: "Task A", Duration: 5 * time.Minute}
	taskB := &Activity{ID: uuid.NewV7(), ProjectID: uuid.NewV7(), DisplayName: "Task B", Duration: 3 * time.Minute}
	t.rows = append(t.rows, NewTableRow(taskA, []*Activity{}))
	t.rows = append(t.rows, NewTableRow(taskB, []*Activity{taskA}))
	return home(t)
}

func (v *View) DisplayDependenciesToAdd(ctx context.Context, query string, projectId uuid.UUID, successorId uuid.UUID) (templ.Component, error) {
	storeActivities, err := v.store.Activity.GetByNameAndProject(ctx, store.GetActivityByNameAndProjectParams{ProjectID: projectId, DisplayName: query})
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			// activity doesn't exist
			return displayDependenciesToAddNotFound(), nil
		}
		return nil, err
	}

	activities := make([]Activity, len(storeActivities))

	for i, storeActivity := range storeActivities {
		activities[i] = Activity{
			ID:          storeActivity.ID,
			ProjectID:   storeActivity.ProjectID,
			DisplayName: storeActivity.DisplayName,
			Duration:    storeActivity.Duration,
		}
	}

	return displayDependenciesToAdd(activities, successorId.String()), nil
}
