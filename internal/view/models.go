package view

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

func NewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *models.ViewActivity {
	return &models.ViewActivity{Id: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}

func NewTable(rows []*models.ViewTableRow) *models.ViewTable {
	return &models.ViewTable{Rows: rows}
}

func NewTableFromStorage(storeActivities []models.StoreActivity, storeDependencies map[models.StoreActivity][]models.StoreActivity) *models.ViewTable {
	storeActivityMap := make(map[models.StoreActivity]*models.ViewActivity)

	for _, storeActivity := range storeActivities {
		//convert activity
		storeActivityMap[storeActivity] = NewActivity(storeActivity.ID, storeActivity.ProjectID, storeActivity.DisplayName, storeActivity.Duration)
	}

	table := models.ViewTable{}

	for _, storeActivity := range storeActivities {
		viewActivity := storeActivityMap[storeActivity]
		storeDependency := storeDependencies[storeActivity]
		viewDependency := make([]*models.ViewActivity, len(storeDependency))
		for i, predecessorActivity := range storeDependency {
			viewDependency[i] = storeActivityMap[predecessorActivity]
		}
		table.Rows = append(table.Rows, NewTableRow(viewActivity, viewDependency))
	}

	return &table
}

func toPostEmptyTableRow(t *models.ViewTableRow) models.PostEmptyTableRow {
	return models.PostEmptyTableRow{ProjectId: t.Activity.ProjectID}
}

func NewTableRow(activity *models.ViewActivity, dependencies []*models.ViewActivity) *models.ViewTableRow {
	return &models.ViewTableRow{Activity: activity, Dependencies: dependencies}
}
