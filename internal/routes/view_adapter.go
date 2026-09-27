package routes

import (
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/view"
)

type httpViewAdapter struct {
	view *view.View
}

// DisplayEmptyTableRow implements [HttpView].
func (h *httpViewAdapter) DisplayEmptyTableRow(params models.ViewDisplayEmptyTableRowParams) {
	h.view.DisplayEmptyTableRow(params)
}

// Home implements [HttpView].
func (h *httpViewAdapter) Home(params models.ViewHomeParams) {
	h.view.Home(params)
}

func NewViewAdapter(view *view.View) *httpViewAdapter {
	return &httpViewAdapter{view: view}
}

var _ HttpView = (*httpViewAdapter)(nil) //ensures viewAdapter implements HttpView at compile time
