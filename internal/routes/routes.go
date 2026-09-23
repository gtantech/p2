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
				log.Printf("returned http internal server error while creating new project. encountered error: %v\n", err)
				return
			}
			storeProjects = []store.Project{storeProject}
		} else {
			http.Error(w, "failed to get projects", http.StatusInternalServerError)
			log.Printf("returned http internal server error while getting projects. encountered error: %v\n", err)
			return
		}
	}
	firstProject := storeProjects[0]
	firstProjectId := firstProject.ID

	//get all activities associated with firstProjectId
	storeActivities, err := rt.store.Activity.GetByProjectID(r.Context(), firstProjectId)
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			storeActivity, err := rt.store.Activity.Create(r.Context(), store.CreateActivityParams{ProjectID: firstProjectId, DisplayName: "", Duration: 0})
			if err != nil {
				http.Error(w, "failed to create new activity", http.StatusInternalServerError)
				log.Printf("returned http internal server error while creating activity. encountered error: %v\n", err)
				return
			}
			storeActivities = []store.Activity{storeActivity}
		} else {
			http.Error(w, "failed to get activities", http.StatusInternalServerError)
			log.Printf("returned http internal server error while getting activities. encountered error: %v\n", err)
			return
		}
	}

	//map activity id to activity
	storeActivitiesMap := make(map[uuid.UUID]store.Activity)
	for _, storeActivity := range storeActivities {
		storeActivitiesMap[storeActivity.ID] = storeActivity
	}

	//map a list of predecessor activities to a successor activity
	storeDependenciesMap := make(map[store.Activity][]store.Activity)
	for _, storeActivity := range storeActivities {
		storeDependencies, err := rt.store.Dependency.GetBySuccessor(r.Context(), storeActivity.ID)
		if err != nil {
			if errors.Is(err, store.ErrDependencyNotFound) {
				storeDependenciesMap[storeActivity] = []store.Activity{}
				continue
			} else {
				http.Error(w, "failed to get dependency", http.StatusInternalServerError)
				log.Printf("returned http internal server error while getting dependency. encountered error: %v\n", err)
				return
			}
		}
		for _, storeDependency := range storeDependencies {
			predecessor := storeActivitiesMap[storeDependency.PredecessorActivityID]
			storeDependenciesMap[storeActivity] = append(storeDependenciesMap[storeActivity], predecessor)
		}
	}

	t := view.NewTableFromStorage(storeActivities, storeDependenciesMap)
	renderTemplComponent(rt.view.Home(t), w, r)
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
