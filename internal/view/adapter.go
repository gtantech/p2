package view

import (
	"github.com/gtantech/p2/internal/models/routemodels"
	"github.com/gtantech/p2/internal/routes"
)

type HttpViewTemplAdapter struct {
}

// RenderHome implements [routes.HttpView].
func (h *HttpViewTemplAdapter) RenderHome(params routemodels.HttpHome) {
	Home(params.Table).Render(params.Request.Context(), params.ResponseWriter)
}

var _ routes.HttpView = (*HttpViewTemplAdapter)(nil) //ensures ExampleStruct implements ExampleInterface at compile time
