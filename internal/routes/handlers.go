package routes

import (
	"net/http"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/template"
	"github.com/gtantech/p2/static"
)

type Routes struct {
	store models.Store
}

func New(store models.Store) *Routes {
	r := &Routes{store: store}
	return r
}

func (rt *Routes) GetHome(w http.ResponseWriter, r *http.Request) {
	mockProjectId := uuid.Max()
	props := rt.store.GetTable(mockProjectId)
	rt.ViewHome(w, r, props)
}

func (rt *Routes) ViewHome(w http.ResponseWriter, r *http.Request, props models.Table) {
	template.Home(props, "Project Planner").Render(r.Context(), w)
}

func (rt *Routes) GetHomeStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}
