package routes

import (
	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/view"
)

type Routes struct {
	store *store.Store
}

func NewRoutes(store *store.Store) *Routes {
	return &Routes{store: store}
}

func (r *Routes) HelloWorldHandler() *templ.ComponentHandler {
	return templ.Handler(view.Home())
}
