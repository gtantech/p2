package view

import (
	json "encoding/json/v2"
	"errors"
	"fmt"
	"uuid"

	"github.com/a-h/templ"
)

type View struct {
}

var ErrActivityNotFound = errors.New("activity not found")

func NewView() *View {
	return &View{}
}

func (v *View) Home(table *Table, homeProjectId uuid.UUID) templ.Component {
	return home(table, homeProjectId)
}

func (v *View) DisplayEmptyTableRow(params DisplayEmptyTableRowParams) templ.Component {
	row := NewTableRow(NewActivity(params.ActivityId, params.ProjectId, "", 0), []*Activity{})
	return displayDependencyTableRow(row, params.ProjectId)
}

func marshalParams(in any) string {
	out, err := json.Marshal(in)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json from params: %v", in))
	}
	return string(out)
}
