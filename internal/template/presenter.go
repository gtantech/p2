package template

import (
	"fmt"
	"net/http"
	"strings"
	"uuid"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models"
)

type TemplPresenter struct {
	model models.Model
}

func todependencyTableRow(tr models.TableRow) dependencyTableRow {
	dr := dependencyTableRow{
		Id:    fmt.Sprintf("activity-row-id-%s", tr.ActivityId),
		Class: "activity-row",
		ActivityTextInputParams: &htmlInput{
			Id:           fmt.Sprintf("activity-input-id-%s", tr.ActivityId),
			Class:        "activity-name",
			Name:         "activity_input",
			Placeholder:  "Activity name",
			Autocomplete: Off,
			Value:        tr.DisplayName,
		},
		DependencyWrapperDivParams: &htmlDiv{
			Id:    fmt.Sprintf("dependency-wrapper-id-%s", tr.ActivityId),
			Class: "dependency-wrapper",
		},
		DependencyTextInputParams: &htmlInput{
			Id:           fmt.Sprintf("dependency-input-id-%s", tr.ActivityId),
			Class:        "dependency-input",
			Name:         "dependency_input",
			Placeholder:  "Activity A, Activity B...",
			Autocomplete: Off,
			Value:        strings.Join(tr.Dependencies, ", "),
		},
		DurationTextInputParams: &htmlInput{
			Id:           fmt.Sprintf("duration-input-id-%s", tr.ActivityId),
			Name:         "duration_input",
			Placeholder:  "e.g. 2h",
			Autocomplete: Off,
			Value:        tr.Duration.String(),
		},
	}
	return dr
}

// DisplayTableRow implements [models.Presenter].
func (t *TemplPresenter) DisplayTableRow(tr models.TableRow) http.Handler {
	dr := todependencyTableRow(tr)
	return templ.Handler(DependencyTableRow(dr))
}

func NewTemplPresenter() *TemplPresenter {
	return &TemplPresenter{}
}

// RegisterModel implements [models.Presenter].
func (t *TemplPresenter) RegisterModel(model models.Model) {
	t.model = model
}

// DisplayHomeHandler implements [models.Presenter].
func (t *TemplPresenter) DisplayHomeHandler(projectId uuid.UUID) http.Handler {
	tableRows, _ := t.model.GetTableRows(projectId)
	dependenciesMap := map[uuid.UUID][]string{}
	for _, tr := range tableRows {
		dependencies, _ := t.model.GetDependencies(tr.ActivityId)
		dependenciesMap[tr.ActivityId] = func() []string {
			dependencyStrSlice := make([]string, len(dependencies))
			for i, dependency := range dependencies {
				dependencyStrSlice[i] = dependency.DisplayName
			}
			return dependencyStrSlice
		}()
	}

	h := homeParams{
		Table: dependencyTable{
			ProjectId: projectId,
			Rows: func() []dependencyTableRow {
				dtr := make([]dependencyTableRow, len(tableRows))
				for i, tr := range tableRows {
					dtr[i] = todependencyTableRow(tr)
				}
				return dtr
			}(),
		},
		TableWrapperDivParams: &htmlDiv{
			Id:    fmt.Sprintf("table-wrapper-id-%s", projectId),
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

var _ models.Presenter = (*TemplPresenter)(nil) //ensures TemplPresenter implements models.Presenter at compile time
