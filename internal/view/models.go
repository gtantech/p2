package view

import (
	"time"
	"uuid"
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

func NewTable() *Table {
	return &Table{rows: []*TableRow{}}
}

type TableRow struct {
	activity     *Activity
	dependencies []*Activity
}

func NewTableRow(activity *Activity, dependencies []*Activity) *TableRow {
	return &TableRow{activity: activity, dependencies: dependencies}
}
