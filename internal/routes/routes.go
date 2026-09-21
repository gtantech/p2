package routes

import (
	"net/http"

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
		view:  view.NewView(),
	}
}

func (rt *Routes) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	rt.view.Home().Render(r.Context(), w)
}
