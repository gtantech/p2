package view

import (
	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models/routemodels"
	"github.com/gtantech/p2/internal/routes"
)

type HttpViewTemplAdapter struct {
}

// RenderDependencyEmptySuggestion implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderDependencyEmptySuggestion(params routemodels.HttpRenderTableRowDependencyEmptySuggestion) {
	HxSwapDivWrapChildren(params.DivTargetSwapId, params.DivHxSwapOob).Render(params.Request.Context(), params.ResponseWriter)
}

// RenderDependencyAddActivitySuggestion implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderDependencyAddActivitySuggestion(params routemodels.HttpRenderTableRowDependencyActivitySuggestion) {
	contents := DependencyTableRowAddActivitiesSuggestion(params.ActivityNames) //.Render(params.Request.Context(), params.ResponseWriter)
	ctx := templ.WithChildren(params.Request.Context(), contents)
	HxSwapDivWrapChildren(params.DivTargetSwapId, params.DivHxSwapOob).Render(ctx, params.ResponseWriter)
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
