package view

import (
	"context"
	json "encoding/json/v2"
	"errors"
	"fmt"
	"io"

	"github.com/gtantech/p2/internal/models"
)

type View struct {
}

var ErrActivityNotFound = errors.New("activity not found")

func NewView() *View {
	return &View{}
}

func (v *View) Home(params models.ViewHomeParams, ctx context.Context, w io.Writer) {
	home(params.Table, params.HomeProjectId).Render(ctx, w)
}

func (v *View) DisplayEmptyTableRow(params models.ViewDisplayEmptyTableRowParams, ctx context.Context, w io.Writer) {
	row := models.NewViewTableRow(models.NewViewActivity(params.ActivityId, params.ProjectId, "", 0), []*models.ViewActivity{})
	displayDependencyTableRow(row, params.ProjectId).Render(ctx, w)
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
