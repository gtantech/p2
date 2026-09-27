package view

import (
	json "encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models"
)

type View struct {
}

var ErrActivityNotFound = errors.New("activity not found")

func NewView() *View {
	return &View{}
}

func (v *View) Home(params models.ViewHomeParams) {
	renderTemplComponent(home(params.Table, params.HomeProjectId), params.HttpResponseWriter, params.HttpRequest)
}

func (v *View) DisplayEmptyTableRow(params models.ViewDisplayEmptyTableRowParams) {
	row := models.NewViewTableRow(models.NewViewActivity(params.ActivityId, params.ProjectId, "", 0), []*models.ViewActivity{})
	renderTemplComponent(displayDependencyTableRow(row, params.ProjectId), params.HttpResponseWriter, params.HttpRequest)
}

func marshalParams(in any) string {
	out, err := json.Marshal(in)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal json from params: %v", in))
	}
	return string(out)
}

func marshalParamsToJsonField(in any) string {
	js := fmt.Sprintf("{\"json\":%s}", marshalParams(in))
	return js
}

func toPostEmptyTableRow(t *models.ViewTableRow) models.RoutesPostEmptyTableRow {
	return models.RoutesPostEmptyTableRow{ProjectId: t.Activity.ProjectID}
}

func renderTemplComponent(component templ.Component, w http.ResponseWriter, r *http.Request) {
	component.Render(r.Context(), w)
}
