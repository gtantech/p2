package routes

import (
	"net/http"
	"uuid"

	"github.com/gtantech/p2/internal/models/routemodels"
)

type HttpView interface {
	RenderHome(params routemodels.HttpHome)
}

type ViewRoutes struct {
	view  HttpView
	store StoreView
}

func NewViewRoutes(view HttpView, store StoreView) *ViewRoutes {
	return &ViewRoutes{view: view, store: store}
}

func (rt *ViewRoutes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	mockProjectId := uuid.Max()
	table, err := rt.store.GetDependencyTableByProjectId(mockProjectId, r.Context())
	if err != nil {
		http.Error(w, "failed to get table", http.StatusInternalServerError)
	}
	rt.view.RenderHome(routemodels.HttpHome{Table: &table, ResponseWriter: w, Request: r})
}
