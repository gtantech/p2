package view

import (
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/store"
)

type Activity struct {
	id          uuid.UUID
	projectID   uuid.UUID
	displayName string
	duration    time.Duration
}

func NewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *Activity {
	return &Activity{id: id, projectID: projectId, displayName: displayName, duration: duration}
}

type Table struct {
	rows []*TableRow
}

func NewTable(rows []*TableRow) *Table {
	return &Table{rows: rows}
}

func NewTableFromStorage(storeActivities []store.Activity, storeDependencies map[store.Activity][]store.Activity) *Table {
	storeActivityMap := make(map[store.Activity]*Activity)

	for _, storeActivity := range storeActivities {
		//convert activity
		storeActivityMap[storeActivity] = NewActivity(storeActivity.ID, storeActivity.ProjectID, storeActivity.DisplayName, storeActivity.Duration)
	}

	table := Table{}

	for _, storeActivity := range storeActivities {
		viewActivity := storeActivityMap[storeActivity]
		storeDependency := storeDependencies[storeActivity]
		viewDependency := make([]*Activity, len(storeDependency))
		for i, predecessorActivity := range storeDependency {
			viewDependency[i] = storeActivityMap[predecessorActivity]
		}
		table.rows = append(table.rows, NewTableRow(viewActivity, viewDependency))
	}

	return &table
}

type TableRow struct {
	activity     *Activity
	dependencies []*Activity
}

func (t *TableRow) ToPostEmptyTableRow() models.PostEmptyTableRow {
	return models.PostEmptyTableRow{ProjectId: t.activity.projectID}
}

func NewTableRow(activity *Activity, dependencies []*Activity) *TableRow {
	return &TableRow{activity: activity, dependencies: dependencies}
}

type DisplayEmptyTableRowParams struct {
	ActivityId uuid.UUID
	ProjectId  uuid.UUID
}
