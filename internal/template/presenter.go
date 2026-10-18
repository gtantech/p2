package template

import (
	"fmt"
	"net/http"
	"strings"
	"uuid"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/models/jsonmodels"
)

var _ models.Presenter = (*TemplPresenter)(nil) //ensures TemplPresenter implements models.Presenter at compile time

type TemplPresenter struct {
	model models.Model
}

// HtmxPostActivityDependencySuggestionAddNewHandlerResp implements [models.Presenter].
func (t *TemplPresenter) HtmxPostActivityDependencySuggestionAddNewHandlerResp(createdTableRows []models.TableRow, rowActivityId uuid.UUID) http.Handler {
	components := []templ.Component{}
	dependencySuggestionDiv := &htmxDiv{
		Id:        fmt.Sprintf("dependency-suggestions-id-%s", rowActivityId),
		Class:     "autocomplete",
		Role:      "listbox",
		HxSwapOob: "innerHTML",
	}

	clearDropdown := DivWrapChildren(dependencySuggestionDiv)
	components = append(components, clearDropdown)

	if len(createdTableRows) > 0 {
		dependencyTableRows := []dependencyTableRow{}
		for _, tr := range createdTableRows {
			dependencyTableRows = append(dependencyTableRows, todependencyTableRow(tr))
		}
		b := TemplateDependencyTableBody(dependencyTableBody{
			ProjectId: createdTableRows[0].ProjectId,
			hxAttr: hx{
				HxSwapOob: fmt.Sprintf("beforeend:#activity-table-body-%s", createdTableRows[0].ProjectId),
			},
			Rows: dependencyTableRows,
		})

		components = append(components, b)
	}

	return templ.Handler(MultiComponent(components))
}

// HtmxOobDependencySuggestionDropdown implements [models.Presenter].
func (t *TemplPresenter) HtmxPostActivityDependencyHandlerResp(suggestedActivityNames []string, suggestedActivityFinalIndex []int, rowActivityId uuid.UUID, projectId uuid.UUID) http.Handler {
	components := []templ.Component{}
	dependencySuggestionDiv := &htmxDiv{
		Id:        fmt.Sprintf("dependency-suggestions-id-%s", rowActivityId),
		Class:     "autocomplete",
		Role:      "listbox",
		HxSwapOob: "innerHTML",
	}
	dropDown := DependencyTableRowDependencyFieldDropdown(
		dependencySuggestionDiv,
		&htmxButton{
			Class:  "suggestion",
			Title:  "Add suggested activities",
			Text:   fmt.Sprintf("Add new activity: %v", strings.Join(suggestedActivityNames, ", ")),
			HxPost: "/form/json/table/row/activity/dependency/suggestion/selected/add-new",
			HxVals: jsonmodels.MarshalParamsToJsonField(jsonmodels.PostFromDependencySuggestionAddNew{
				ActivityId:                   rowActivityId,
				ProjectId:                    projectId,
				ActivitiesToCreate:           suggestedActivityNames,
				ActivitiesToCreateFinalIndex: suggestedActivityFinalIndex,
			}),
		},
	)
	clearDropdown := DivWrapChildren(dependencySuggestionDiv)
	if len(suggestedActivityNames) > 0 {
		components = append(components, dropDown)
	} else {
		components = append(components, clearDropdown)
	}
	return templ.Handler(MultiComponent(components))
}

// HtmlTableRow implements [models.Presenter].
func (t *TemplPresenter) HtmlTableRow(tr models.TableRow) http.Handler {
	dr := todependencyTableRow(tr)
	return templ.Handler(DependencyTableRow(dr))
}

func NewTemplPresenter() *TemplPresenter {
	return &TemplPresenter{}
}

// HtmlHomeHandler implements [models.Presenter].
func (t *TemplPresenter) HtmlHomeHandler(table models.Table) http.Handler {
	h := homeParams{
		Table: dependencyTable{
			ProjectId: table.ProjectId,
			Rows: func() []dependencyTableRow {
				dtr := make([]dependencyTableRow, len(table.Rows))
				for i, tr := range table.Rows {
					dtr[i] = todependencyTableRow(tr)
				}
				return dtr
			}(),
		},
		TableWrapperDivParams: &htmlDiv{
			Id:    fmt.Sprintf("table-wrapper-id-%s", table.ProjectId),
			Class: "table-wrapper",
		},
		ContainerDivParams: &htmlDiv{
			Class: "container",
		},
		PageTitle: "Project Planner",
	}

	component := Home(h)
	return templ.Handler(component)
}

func todependencyTableRow(tr models.TableRow) dependencyTableRow {
	activityNameDomName := "activity_input"
	activityDurationDomName := "duration_input"
	activityDependencyDomName := "dependency_input"
	dr := dependencyTableRow{
		Id:    fmt.Sprintf("activity-row-id-%s", tr.ActivityId),
		Class: "activity-row",
		ActivityTextInputParams: &htmxInput{
			Id:           fmt.Sprintf("activity-input-id-%s", tr.ActivityId),
			Class:        "activity-name",
			Name:         activityNameDomName,
			Placeholder:  "Activity name",
			Autocomplete: Off,
			Value:        tr.DisplayName,
			HxPut:        "/form/json/table/row/activity/name",
			HxVals:       jsonmodels.MarshalParamsToJsonField(jsonmodels.PutFromRowActivityNameChange{ActivityId: tr.ActivityId, DomName: activityNameDomName}),
			HxTrigger:    "input changed delay:200ms",
		},
		DependencyWrapperDivParams: &htmlDiv{
			Id:    fmt.Sprintf("dependency-wrapper-id-%s", tr.ActivityId),
			Class: "dependency-wrapper",
		},
		DependencyTextInputParams: &htmxInput{
			Id:           fmt.Sprintf("dependency-input-id-%s", tr.ActivityId),
			Class:        "dependency-input",
			Name:         activityDependencyDomName,
			Placeholder:  "Activity A, Activity B...",
			Autocomplete: Off,
			Value:        strings.Join(tr.Dependencies, ", "),
			HxPost:       "/form/json/table/row/activity/dependency",
			HxVals:       jsonmodels.MarshalParamsToJsonField(jsonmodels.PutFromRowActivityDependencyChange{ActivityId: tr.ActivityId, ProjectId: tr.ProjectId, DomName: activityDependencyDomName}),
			HxTrigger:    "input changed delay:200ms",
		},
		DependencySuggestionDivParams: &htmxDiv{
			Id:    fmt.Sprintf("dependency-suggestions-id-%s", tr.ActivityId),
			Class: "autocomplete",
			Role:  "listbox",
		},
		DurationTextInputParams: &htmxInput{
			Id:           fmt.Sprintf("duration-input-id-%s", tr.ActivityId),
			Name:         activityDurationDomName,
			Placeholder:  "e.g. 2h",
			Autocomplete: Off,
			Value:        tr.Duration.String(),
			HxPut:        "/form/json/table/row/activity/duration",
			HxVals:       jsonmodels.MarshalParamsToJsonField(jsonmodels.PutFromRowActivityDurationChange{ActivityId: tr.ActivityId, DomName: activityDurationDomName}),
			HxTrigger:    "input changed delay:200ms",
		},
		TableRowAddBtnParams: &htmxButton{
			Class:    "add-button",
			Title:    "Add activity",
			Text:     "+",
			HxPost:   "/form/json/table/row/empty/below",
			HxVals:   jsonmodels.MarshalParamsToJsonField(jsonmodels.PostFromRowPlusBtn{ProjectId: tr.ProjectId, RelativeToActivityId: tr.ActivityId}),
			HxTarget: "closest tr",
			HxSwap:   "afterend",
		},
	}
	return dr
}
