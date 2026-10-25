package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/template"
	"github.com/gtantech/p2/static"
)

type Routes struct {
	store Store
}

func New(store Store) *Routes {
	r := &Routes{store: store}
	return r
}

func (rt *Routes) GetHome(w http.ResponseWriter, r *http.Request) {
	homeProject, err := rt.store.GetHomeProject(r.Context())
	if err != nil {
		if storeErrorHandler(err, w, "failed to get home project in GetHome handler") {
			return
		}
	}
	props, err := rt.store.GetTable(r.Context(), homeProject.ProjecId)
	if err != nil {
		if storeErrorHandler(err, w, "failed to get home table in GetHome handler") {
			return
		}
	}
	rt.ViewHome(w, r, props)
}

func (rt *Routes) ViewHome(w http.ResponseWriter, r *http.Request, props models.Table) {
	template.Home(props, "Project Planner").Render(r.Context(), w)
}

func (rt *Routes) GetHomeStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}
