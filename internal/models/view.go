package models

import (
	"time"
	"uuid"
)

type ViewActivity struct {
	Id          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type ViewTable struct {
	Rows      []*ViewTableRow
	ProjectId uuid.UUID
}

type ViewTableRow struct {
	Activity     *ViewActivity
	Dependencies []*ViewActivity
}

type ViewDisplayEmptyTableRowParams struct {
	ActivityId uuid.UUID
	ProjectId  uuid.UUID
}

type ViewHomeParams struct {
	Table         *ViewTable
	HomeProjectId uuid.UUID
}

type ViewDisplayDependencyWrapper struct {
	RowActivityId   uuid.UUID
	ProjectId       uuid.UUID
	DependencyNames []string
}

type ViewDisplayDependencyWrapperWithNewActivitySuggestion struct {
	ViewDisplayDependencyWrapper
	NewActivityDispName string
}

func NewViewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *ViewActivity {
	return &ViewActivity{Id: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}

func NewViewTable(rows []*ViewTableRow) *ViewTable {
	return &ViewTable{Rows: rows}
}

func NewViewTableFromStorage(storeActivities []StoreActivity, storeDependencies map[StoreActivity][]StoreActivity) *ViewTable {
	storeActivityMap := make(map[StoreActivity]*ViewActivity)

	for _, storeActivity := range storeActivities {
		//convert activity
		storeActivityMap[storeActivity] = NewViewActivity(storeActivity.ID, storeActivity.ProjectID, storeActivity.DisplayName, storeActivity.Duration)
	}

	table := ViewTable{}

	for _, storeActivity := range storeActivities {
		viewActivity := storeActivityMap[storeActivity]
		storeDependency := storeDependencies[storeActivity]
		viewDependency := make([]*ViewActivity, len(storeDependency))
		for i, predecessorActivity := range storeDependency {
			viewDependency[i] = storeActivityMap[predecessorActivity]
		}
		table.Rows = append(table.Rows, NewViewTableRow(viewActivity, viewDependency))
	}

	return &table
}

func NewViewTableRow(activity *ViewActivity, dependencies []*ViewActivity) *ViewTableRow {
	return &ViewTableRow{Activity: activity, Dependencies: dependencies}
}
