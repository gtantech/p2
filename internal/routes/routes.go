package routes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/view"
)

type Routes struct {
	store *store.Store
	view  *view.View
}

func NewRoutes(store *store.Store) *Routes {
	return &Routes{
		store: store,
		view:  view.NewView(store),
	}
}

func renderTemplComponent(component templ.Component, w http.ResponseWriter, r *http.Request) {
	component.Render(r.Context(), w)
}

func (rt *Routes) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplComponent(rt.view.Home(), w, r)
}
