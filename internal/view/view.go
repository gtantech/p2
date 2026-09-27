package view

import (
	json "encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"time"
	"uuid"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models"
)

type View struct {
}

var ErrActivityNotFound = errors.New("activity not found")

func NewView() *View {
	return &View{}
}

func (v *View) Home(params models.ViewHomeParams) {
	renderTemplComponent(home(params.Table, params.HomeProjectId), params.HttpResponseWriter, params.HttpRequest)
}

func (v *View) DisplayEmptyTableRow(params models.ViewDisplayEmptyTableRowParams) {
	row := NewViewTableRow(NewViewActivity(params.ActivityId, params.ProjectId, "", 0), []*models.ViewActivity{})
	renderTemplComponent(displayDependencyTableRow(row, params.ProjectId), params.HttpResponseWriter, params.HttpRequest)
}

func marshalParams(in any) string {
	out, err := json.Marshal(in)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json from params: %v", in))
	}
	return string(out)
}

func NewViewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *models.ViewActivity {
	return &models.ViewActivity{Id: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}

func NewViewTable(rows []*models.ViewTableRow) *models.ViewTable {
	return &models.ViewTable{Rows: rows}
}

func NewViewTableFromStorage(storeActivities []models.StoreActivity, storeDependencies map[models.StoreActivity][]models.StoreActivity) *models.ViewTable {
	storeActivityMap := make(map[models.StoreActivity]*models.ViewActivity)

	for _, storeActivity := range storeActivities {
		//convert activity
		storeActivityMap[storeActivity] = NewViewActivity(storeActivity.ID, storeActivity.ProjectID, storeActivity.DisplayName, storeActivity.Duration)
	}

	table := models.ViewTable{}

	for _, storeActivity := range storeActivities {
		viewActivity := storeActivityMap[storeActivity]
		storeDependency := storeDependencies[storeActivity]
		viewDependency := make([]*models.ViewActivity, len(storeDependency))
		for i, predecessorActivity := range storeDependency {
			viewDependency[i] = storeActivityMap[predecessorActivity]
		}
		table.Rows = append(table.Rows, NewViewTableRow(viewActivity, viewDependency))
	}

	return &table
}

func toPostEmptyTableRow(t *models.ViewTableRow) models.RoutesPostEmptyTableRow {
	return models.RoutesPostEmptyTableRow{ProjectId: t.Activity.ProjectID}
}

func NewViewTableRow(activity *models.ViewActivity, dependencies []*models.ViewActivity) *models.ViewTableRow {
	return &models.ViewTableRow{Activity: activity, Dependencies: dependencies}
}

func renderTemplComponent(component templ.Component, w http.ResponseWriter, r *http.Request) {
	component.Render(r.Context(), w)
}
