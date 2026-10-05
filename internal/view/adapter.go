package view

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/models/routemodels"
	"github.com/gtantech/p2/internal/routes"
)

type HttpViewTemplAdapter struct {
}

// RenderDependencyInputSelectedAddActivityResp implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderDependencyInputSelectedAddActivityResp(params routemodels.HttpRenderTableRowDependencyAddActivitySelectedResp) {
	components := []templ.Component{}

	components = append(components, HxSwapDivWrapChildren(fmt.Sprintf("dependency-suggestions-%s", params.FromRowActivityId), "innerHTML"))

	MultiComponent(components).Render(params.Request.Context(), params.ResponseWriter)
}

// RenderDependencyAddActivitySuggestion implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderDependencyInputResp(params routemodels.HttpRenderTableRowDependencyInputResp) {
	components := []templ.Component{}

	if len(params.ActivityNames) > 0 {
		components = append(components, DependencyTableRowAddActivitiesSuggestion(params.ProjectId, params.RowActivityId, params.ActivityNames, params.DivTargetSwapId, params.DivHxSwapOob))
	} else {
		components = append(components, HxSwapDivWrapChildren(params.DivTargetSwapId, params.DivHxSwapOob))
	}

	MultiComponent(components).Render(params.Request.Context(), params.ResponseWriter)
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
