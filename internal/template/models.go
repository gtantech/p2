package template

import (
	"uuid"
)

type InputAutocomplete string

const (
	On  InputAutocomplete = "on"
	Off InputAutocomplete = "off"
)

type htmlInput struct {
	Id           string
	Name         string
	Class        string
	Placeholder  string
	Autocomplete InputAutocomplete
	Value        string
}

type htmlDiv struct {
	Id    string
	Class string
}

type htmlButton struct {
	Class string
	Title string
	Text  string
}

type dependencyTableRow struct {
	Id                         string
	Class                      string
	ActivityTextInputParams    *htmlInput
	DependencyWrapperDivParams *htmlDiv
	DependencyTextInputParams  *htmlInput
	DurationTextInputParams    *htmlInput
	TableRowAddBtnParams       *htmlButton
}

type dependencyTable struct {
	ProjectId uuid.UUID
	Rows      []dependencyTableRow
}

type homeParams struct {
	Table                 dependencyTable
	TableWrapperDivParams *htmlDiv
	ContainerDivParams    *htmlDiv
	PageTitle             string
}
