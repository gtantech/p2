package view

type table struct {
	rows []*tableRow
}

func newTable() *table {
	return &table{rows: []*tableRow{}}
}

type tableRow struct {
	activity     *activity
	dependencies []*activity
}

func newTableRow(activity *activity, dependencies []*activity) *tableRow {
	return &tableRow{activity: activity, dependencies: dependencies}
}
