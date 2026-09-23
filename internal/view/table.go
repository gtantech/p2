package view

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
