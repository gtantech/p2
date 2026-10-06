package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/models"
)

type Routes struct {
	presenter models.Presenter
	store     models.Model
}

func NewRoutes(presenter models.Presenter, store models.Model) *Routes {
	r := &Routes{presenter: presenter, store: store}
	presenter.RegisterModel(store)
	return r
}

func (rt *Routes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	rt.presenter.DisplayHomeHandler().ServeHTTP(w, r)
}
