package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/models/routemodels"
)

type HttpView interface {
	RenderHome(params routemodels.HttpHome)
}

type ViewRoutes struct {
	view HttpView
}

func NewViewRoutes(view HttpView) *ViewRoutes {
	return &ViewRoutes{view: view}
}

func (rt *ViewRoutes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	rt.view.RenderHome(routemodels.HttpHome{ResponseWriter: w, Request: r})
}
