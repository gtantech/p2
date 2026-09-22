package routes

import (
	"errors"
	"log"
	"net/http"
	"uuid"

	"github.com/a-h/templ"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/view"
	"github.com/gtantech/p2/static"
)

type Routes struct {
	store *store.Store
	view  *view.View
}

func NewRoutes(store *store.Store) *Routes {
	return &Routes{
		store: store,
		view:  view.NewView(store),
	}
}

func renderTemplComponent(component templ.Component, w http.ResponseWriter, r *http.Request) {
	component.Render(r.Context(), w)
}

func (rt *Routes) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplComponent(rt.view.Home(), w, r)
}

func (rt *Routes) DisplayDependenciesToAdd(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		http.Error(w, "missing search parameter", http.StatusBadRequest)
		return
	}
	projectId, err := uuid.Parse(r.URL.Query().Get("project-id"))
	if err != nil {
		http.Error(w, "invalid project id parameter", http.StatusBadRequest)
		return
	}
	successorId, err := uuid.Parse(r.URL.Query().Get("successor-id"))

	if err != nil {
		http.Error(w, "invalid successor id parameter", http.StatusBadRequest)
		return
	}
	component, err := rt.view.DisplayDependenciesToAdd(r.Context(), search, projectId, successorId)
	if err != nil {
		if errors.Is(err, view.ErrActivityNotFound) {
			http.Error(w, "activity not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("internal server error occurred with error:%v\n", err)
		return
	}
	renderTemplComponent(component, w, r)
}

func (rt *Routes) GetHomeStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}
