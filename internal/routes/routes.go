package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/view"
)

type Routes struct {
	store *store.Store
}

func NewRoutes(store *store.Store) *Routes {
	return &Routes{store: store}
}

func (rt *Routes) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	view.Home().Render(r.Context(), w)
}
