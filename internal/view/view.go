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

func (v *View) Home() templ.Component {
	return home()
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

	activities := make([]activity, len(storeActivities))

	for i, storeActivity := range storeActivities {
		activities[i] = activity{
			ID:          storeActivity.ID,
			ProjectID:   storeActivity.ProjectID,
			DisplayName: storeActivity.DisplayName,
			Duration:    storeActivity.Duration,
		}
	}

	return displayDependenciesToAdd(activities, successorId.String()), nil
}
