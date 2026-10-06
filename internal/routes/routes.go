package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/models"
)

type Routes struct {
	presenter models.Presenter
	models.Model
}

func NewRoutes(presenter models.Presenter, store models.Model) *Routes {
	r := &Routes{presenter: presenter, Model: store}
	presenter.RegisterModel(r)
	return r
}

func (rt *Routes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	rt.presenter.DisplayHomeHandler().ServeHTTP(w, r)
}

var _ models.Model = (*Routes)(nil) //ensures Routes implements models.Model at compile time
