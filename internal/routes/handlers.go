package routes

import (
	"net/http"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/static"
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
	homeProjectId := uuid.Max() //mock home project id
	rt.presenter.DisplayHomeHandler(homeProjectId).ServeHTTP(w, r)
}

func (rt *Routes) GetHomeStyleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}
