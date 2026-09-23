package view

import (
	"context"
	"errors"
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

func (v *View) Home(table *Table) templ.Component {
	return home(table)
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
			id:          storeActivity.ID,
			projectID:   storeActivity.ProjectID,
			displayName: storeActivity.DisplayName,
			duration:    storeActivity.Duration,
		}
	}

	return displayDependenciesToAdd(activities, successorId.String()), nil
}
