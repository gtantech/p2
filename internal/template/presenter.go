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

func NewTemplPresenter() *TemplPresenter {
	return &TemplPresenter{}
}

// RegisterModel implements [models.Presenter].
func (t *TemplPresenter) RegisterModel(model models.Model) {
	t.model = model
}

// DisplayHomeHandler implements [models.Presenter].
func (t *TemplPresenter) DisplayHomeHandler(projectId uuid.UUID) http.Handler {
	activities, _ := t.model.GetActivities(projectId)
	dependenciesMap := map[uuid.UUID][]string{}
	for _, activity := range activities {
		dependencies, _ := t.model.GetDependencies(activity.ActivityId)
		dependenciesMap[activity.ActivityId] = func() []string {
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
				dtr := make([]dependencyTableRow, len(activities))
				for i, activity := range activities {
					dtr[i] = dependencyTableRow{
						Id:    fmt.Sprintf("activity-row-id-%s", activity.ActivityId),
						Class: "activity-row",
						ActivityTextInputParams: &htmlInput{
							Id:           fmt.Sprintf("activity-input-id-%s", activity.ActivityId),
							Class:        "activity-name",
							Name:         "activity_input",
							Placeholder:  "Activity name",
							Autocomplete: Off,
							Value:        activity.DisplayName,
						},
						DependencyWrapperDivParams: &htmlDiv{
							Id:    fmt.Sprintf("dependency-wrapper-id-%s", activity.ActivityId),
							Class: "dependency-wrapper",
						},
						DependencyTextInputParams: &htmlInput{
							Id:           fmt.Sprintf("dependency-input-id-%s", activity.ActivityId),
							Class:        "dependency-input",
							Name:         "dependency_input",
							Placeholder:  "Activity A, Activity B...",
							Autocomplete: Off,
							Value:        strings.Join(dependenciesMap[activity.ActivityId], ", "),
						},
						DurationTextInputParams: &htmlInput{
							Id:           fmt.Sprintf("duration-input-id-%s", activity.ActivityId),
							Name:         "duration_input",
							Placeholder:  "e.g. 2h",
							Autocomplete: Off,
							Value:        activity.Duration.String(),
						},
					}
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
