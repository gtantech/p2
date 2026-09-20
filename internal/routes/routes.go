package routes

import (
	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/view"
)

type Routes struct {
	store store.Store
}

func (r *Routes) HelloWorldHandler() *templ.ComponentHandler {
	return templ.Handler(view.Home())
}
