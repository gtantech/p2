package routes

import (
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/view"
)

type httpViewAdapter struct {
	view *view.View
}

// DisplayDependencyWrapper implements [HttpView].
func (h *httpViewAdapter) DisplayDependencyWrapper(params models.HttpViewDisplayDependencyWrapper) {
	view.NewView().DisplayDependencyWrapper(
		models.ViewDisplayDependencyWrapper{
			RowActivityId:   params.RowActivityId,
			ProjectId:       params.ProjectId,
			DependencyNames: params.DependencyNames,
		},
		params.HttpRequest.Context(),
		params.HttpResponseWriter,
	)
}

// DisplayDependencyWithNewActivitySuggestion implements [HttpView].
func (h *httpViewAdapter) DisplayDependencyWrapperWithNewActivitySuggestion(params models.HttpViewDisplayDependencyWrapperWithNewActivitySuggestion) {
	h.view.DisplayDependencyWrapperWithNewActivitySuggestion(
		models.ViewDisplayDependencyWrapperWithNewActivitySuggestion{
			RowActivityId:       params.RowActivityId,
			ProjectId:           params.ProjectId,
			DependencyNames:     params.DependencyNames,
			NewActivityDispName: params.NewActivityDispName,
		},
		params.HttpRequest.Context(),
		params.HttpResponseWriter,
	)
}

// DisplayEmptyTableRow implements [HttpView].
func (h *httpViewAdapter) DisplayEmptyTableRow(params models.HttpViewDisplayEmptyTableRowParams) {
	h.view.DisplayEmptyTableRow(models.ViewDisplayEmptyTableRowParams{ActivityId: params.ActivityId, ProjectId: params.ProjectId}, params.HttpRequest.Context(), params.HttpResponseWriter)
}

// DisplayHome implements [HttpView].
func (h *httpViewAdapter) DisplayHome(params models.HttpViewHomeParams) {
	h.view.Home(models.ViewHomeParams{Table: params.Table}, params.HttpRequest.Context(), params.HttpResponseWriter)
}

func NewViewAdapter(view *view.View) *httpViewAdapter {
	return &httpViewAdapter{view: view}
}

var _ HttpView = (*httpViewAdapter)(nil) //ensures viewAdapter implements HttpView at compile time
