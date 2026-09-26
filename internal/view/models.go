package view

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

func NewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *models.Activity {
	return &models.Activity{Id: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}

func NewTable(rows []*models.TableRow) *models.Table {
	return &models.Table{Rows: rows}
}

func NewTableFromStorage(storeActivities []models.StoreActivity, storeDependencies map[models.StoreActivity][]models.StoreActivity) *models.Table {
	storeActivityMap := make(map[models.StoreActivity]*models.Activity)

	for _, storeActivity := range storeActivities {
		//convert activity
		storeActivityMap[storeActivity] = NewActivity(storeActivity.ID, storeActivity.ProjectID, storeActivity.DisplayName, storeActivity.Duration)
	}

	table := models.Table{}

	for _, storeActivity := range storeActivities {
		viewActivity := storeActivityMap[storeActivity]
		storeDependency := storeDependencies[storeActivity]
		viewDependency := make([]*models.Activity, len(storeDependency))
		for i, predecessorActivity := range storeDependency {
			viewDependency[i] = storeActivityMap[predecessorActivity]
		}
		table.Rows = append(table.Rows, NewTableRow(viewActivity, viewDependency))
	}

	return &table
}

func toPostEmptyTableRow(t *models.TableRow) models.PostEmptyTableRow {
	return models.PostEmptyTableRow{ProjectId: t.Activity.ProjectID}
}

func NewTableRow(activity *models.Activity, dependencies []*models.Activity) *models.TableRow {
	return &models.TableRow{Activity: activity, Dependencies: dependencies}
}

type DisplayEmptyTableRowParams struct {
	ActivityId uuid.UUID
	ProjectId  uuid.UUID
}
