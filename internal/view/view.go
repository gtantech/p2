package view

import (
	json "encoding/json/v2"
	"errors"
	"fmt"
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
