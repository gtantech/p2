package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/template"
)

type Routes struct {
	store models.Model
}

func NewRoutes(store models.Model) *Routes {
	r := &Routes{store: store}
	return r
}

func (rt *Routes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	rt.ViewHome(w, r, models.Home{Name: "World"})
}

func (rt *Routes) ViewHome(w http.ResponseWriter, r *http.Request, props models.Home) {
	template.Hello(props.Name).Render(r.Context(), w)
}
