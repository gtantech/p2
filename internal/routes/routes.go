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
	storeProjects, err := rt.store.Project.GetProjects(r.Context())
	if err != nil {
		if errors.Is(err, store.ErrProjectNotFound) {
			storeProject, err := rt.store.Project.Create(r.Context(), store.CreateProjectParams{DisplayName: "Project 1"})
			if err != nil {
				http.Error(w, "failed to create new project", http.StatusInternalServerError)
				return
			}
			storeProjects = []store.Project{storeProject}
		} else {
			http.Error(w, "failed to get projects", http.StatusInternalServerError)
			return
		}
	}
	renderTemplComponent(rt.view.Home(storeProjects[0].ID), w, r)
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
