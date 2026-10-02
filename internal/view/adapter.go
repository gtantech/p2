package view

import (
	"github.com/gtantech/p2/internal/models/routemodels"
	"github.com/gtantech/p2/internal/routes"
)

type HttpViewTemplAdapter struct {
}

// RenderDependencyAddActivitySuggestion implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderDependencyAddActivitySuggestion(params routemodels.HttpRenderTableRowDependencyActivitySuggestion) {
	DependencyTableRowAddActivitiesSuggestion(params.ActivityNames, params.RowActivityId, params.HxTarget, params.HxSwap).Render(params.Request.Context(), params.ResponseWriter)
}

// RenderDependencyTableRow implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderDependencyTableRow(params routemodels.HttpRenderTableRow) {
	DependencyTableRow(&params.TableRow, params.NextRowSortRank).Render(params.Request.Context(), params.ResponseWriter)
}

// RenderHome implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderHome(params routemodels.HttpHome) {
	Home(params.Table).Render(params.Request.Context(), params.ResponseWriter)
}

var _ routes.HttpView = (*HttpViewTemplAdapter)(nil) //ensures ExampleStruct implements ExampleInterface at compile time
