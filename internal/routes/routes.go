package routes

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/store"
	"github.com/gtantech/p2/internal/view"
	"github.com/gtantech/p2/static"
)

type Routes struct {
	store StoreService
	view  *view.View
}

func NewRoutes(store StoreService) *Routes {
	return &Routes{
		store: store,
		view:  view.NewView(),
	}
}

func renderTemplComponent(component templ.Component, w http.ResponseWriter, r *http.Request) {
	component.Render(r.Context(), w)
}

func (rt *Routes) PutActivityDependencyUpdateFromTableHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	activityId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "malformed id", http.StatusBadRequest)
		return
	}

	relationship := models.FS
	if relationship == "" {
		http.Error(w, "missing relationship parameter", http.StatusBadRequest)
		return
	}

	storeActivity, err := rt.store.Activity().GetByID(r.Context(), activityId)
	if err != nil {
		http.Error(w, "failed to get activity", http.StatusInternalServerError)
	}

	storeDependencies, err := rt.store.Dependency().GetPredecessorNamesBySuccessor(r.Context(), activityId)
	if err != nil {
		if !errors.Is(err, store.ErrDependencyNotFound) {
			http.Error(w, "failed to get dependencies", http.StatusInternalServerError)
			log.Printf("returned http internal server error while getting projects. encountered error: %v\n", err)
			return
		}
	}

	predecessorNameToDependencyId := make(map[string]uuid.UUID)

	for _, dependency := range storeDependencies {
		predecessorNameToDependencyId[dependency.PredecessorActivityName] = dependency.DependencyID
	}

	dependency_input := r.FormValue("dependency_input")
	if dependency_input == "" {
		//remove all dependencies
		for _, dependency := range storeDependencies {
			rt.store.Dependency().Delete(r.Context(), dependency.DependencyID)
		}
		w.Write([]byte("OK"))
		return
	}
	parts := strings.Split(dependency_input, ",")

	dependencyInputMap := make(map[string]bool)
	for _, part := range parts {
		dependencyInputMap[strings.TrimSpace(part)] = true
	}

	// check if user deleted value
	for key := range predecessorNameToDependencyId {
		//if key in predecessorNameToDependencyId is not in input, user has deleted value
		if _, ok := dependencyInputMap[key]; !ok {
			deleteId := predecessorNameToDependencyId[key]
			rt.store.Dependency().Delete(r.Context(), deleteId)
		}
	}

	// check if user added new value
	for key := range dependencyInputMap {
		//if key in dependencyInputMap is not in predecessorNameToDependencyId, user has added value
		if _, ok := predecessorNameToDependencyId[key]; !ok {
			findUserSpecifiedActivity, err := rt.store.Activity().GetByNameAndProject(r.Context(), store.GetActivityByNameAndProjectParams{
				ProjectID:   storeActivity.ProjectID,
				DisplayName: key})
			if err != nil {
				http.Error(w, "failed to get activity", http.StatusBadRequest)
				return
			}
			rt.store.Dependency().Create(r.Context(), store.CreateDepdencencyParams{
				ProjectID:             storeActivity.ProjectID,
				Relationship:          models.RelationshipType(relationship),
				PredecessorActivityID: findUserSpecifiedActivity[0].ID,
				SuccessorActivityID:   activityId,
			})
		}
	}
}

func (rt *Routes) PutActivityDurationUpdateFromTableHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	activityId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "malformed id", http.StatusBadRequest)
		return
	}

	duration_input := r.FormValue("duration_input")
	if duration_input == "" {
		_, err = rt.store.Activity().Update(r.Context(), store.UpdateActivityParams{Id: activityId, DisplayName: duration_input, Duration: 0})
		if err != nil {
			http.Error(w, "failed to update activity", http.StatusInternalServerError)
			log.Printf("returned http internal server error while updating activity %v. encountered error: %v\n", activityId, err)
		}
		return
	}
	duration, err := time.ParseDuration(duration_input)
	if err != nil {
		http.Error(w, "invalid duration input", http.StatusBadRequest)
	}
	storeActivity, err := rt.store.Activity().GetByID(r.Context(), activityId)

	if err != nil {
		http.Error(w, "failed to get activity", http.StatusInternalServerError)
	}
	_, err = rt.store.Activity().Update(r.Context(), store.UpdateActivityParams{Id: activityId, DisplayName: storeActivity.DisplayName, Duration: duration})
	if err != nil {
		http.Error(w, "failed to update activity", http.StatusInternalServerError)
	}
}

func (rt *Routes) PutActivityNameUpdateFromTableHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	activityId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "malformed id", http.StatusBadRequest)
		return
	}

	name := r.FormValue("activity_input")
	if name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	storeActivity, err := rt.store.Activity().GetByID(r.Context(), activityId)

	if err != nil {
		http.Error(w, "failed to get activity", http.StatusInternalServerError)
	}
	_, err = rt.store.Activity().Update(r.Context(), store.UpdateActivityParams{Id: activityId, DisplayName: name, Duration: storeActivity.Duration})
	if err != nil {
		http.Error(w, "failed to update activity", http.StatusInternalServerError)
	}
}

func (rt *Routes) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	storeProjects, err := rt.store.Project().GetProjects(r.Context())
	if err != nil {
		if errors.Is(err, store.ErrProjectNotFound) {
			storeProject, err := rt.store.Project().Create(r.Context(), store.CreateProjectParams{DisplayName: "Project 1"})
			if err != nil {
				http.Error(w, "failed to create new project", http.StatusInternalServerError)
				log.Printf("returned http internal server error while creating new project. encountered error: %v\n", err)
				return
			}
			storeProjects = []models.StoreProject{storeProject}
		} else {
			http.Error(w, "failed to get projects", http.StatusInternalServerError)
			log.Printf("returned http internal server error while getting projects. encountered error: %v\n", err)
			return
		}
	}
	firstProject := storeProjects[0]
	firstProjectId := firstProject.ID

	//get all activities associated with firstProjectId
	storeActivities, err := rt.store.Activity().GetByProjectID(r.Context(), firstProjectId)
	if err != nil {
		if errors.Is(err, store.ErrActivityNotFound) {
			storeActivity, err := rt.store.Activity().Create(r.Context(), store.CreateActivityParams{ProjectID: firstProjectId, DisplayName: "", Duration: 0})
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
		storeDependencies, err := rt.store.Dependency().GetBySuccessor(r.Context(), storeActivity.ID)
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
	renderTemplComponent(rt.view.Home(t, firstProjectId), w, r)
}

func (rt *Routes) GetHomeStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}

func (rt *Routes) PostEmptyTableRow(w http.ResponseWriter, r *http.Request) {
	projectId, err := uuid.Parse(r.FormValue("projectId"))
	if err != nil {
		http.Error(w, "invalid project id parameter", http.StatusBadRequest)
		return
	}

	storeActivity, err := rt.store.Activity().Create(r.Context(), store.CreateActivityParams{ProjectID: projectId, DisplayName: "", Duration: 0})

	if err != nil {
		//TODO check for duplicate name error
		http.Error(w, "database returned error", http.StatusInternalServerError)
		return
	}

	renderTemplComponent(rt.view.DisplayEmptyTableRow(view.DisplayEmptyTableRowParams{ActivityId: storeActivity.ID, ProjectId: storeActivity.ProjectID}), w, r)
}
