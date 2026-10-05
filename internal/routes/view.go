package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"uuid"

	"github.com/gtantech/p2/internal/models/jsonmodels"
	"github.com/gtantech/p2/internal/models/routemodels"
	"github.com/gtantech/p2/internal/models/storemodels"
	"github.com/gtantech/p2/internal/models/viewmodels"
)

type HttpView interface {
	RenderHome(params routemodels.HttpHome)
	RenderDependencyTableRow(params routemodels.HttpRenderTableRow)
	RenderDependencyInputResp(params routemodels.HttpRenderTableRowDependencyInputResp)
	RenderDependencyInputSelectedAddActivityResp(params routemodels.HttpRenderTableRowDependencyAddActivitySelectedResp)
}

type ViewRoutes struct {
	view      HttpView
	storeView StoreView
	store     Store
}

func NewViewRoutes(view HttpView, storeView StoreView, store Store) *ViewRoutes {
	return &ViewRoutes{view: view, storeView: storeView, store: store}
}

func (rt *ViewRoutes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	mockProjectId := uuid.Max()
	table, err := rt.storeView.GetDependencyTableByProjectId(mockProjectId, r.Context())
	if err != nil {
		http.Error(w, "failed to get table", http.StatusInternalServerError)
	}
	rt.view.RenderHome(routemodels.HttpHome{Table: &table, ResponseWriter: w, Request: r})
}

func (rt *ViewRoutes) PostTableOfDependencyRowDependencyHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	var dto jsonmodels.UpdateDependencyFromInput
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
		return
	}

	//parse user input
	userInputs := strings.Split(r.FormValue(dto.DomName), ",")

	//clean user input
	for i, userInput := range userInputs {
		userInputs[i] = strings.TrimSpace(userInput)
	}

	userInputs = removeDuplicates(userInputs)
	userInputs = removeEmptyString(userInputs)

	//reconcile user input with storeActivities:
	// - by returning a suggestion for creating unknown activities,
	// - deleting activities from depenendency store where no longer in user input
	storeActivities, err := rt.store.GetActivitiesByProjectId(dto.ProjectId, r.Context())
	if err != nil {
		http.Error(w, "failed to get activities from store", http.StatusInternalServerError)
		log.Printf("returned http internal server error while getting activities from store. Encountered error: %v", err)
		return
	}
	storeActivitiesMap := make(map[string]storemodels.Activity)
	for _, storeActivity := range storeActivities {
		storeActivitiesMap[storeActivity.DispName] = storeActivity
	}

	//determine new values (unknown activities)
	userInputsNewValues := valuesNotInMap(userInputs, storeActivitiesMap)

	//delete activities from depenendency store where no longer in user input
	// -by getting activities in store
	// -by checking activities in store that are missing from user input
	// -by deleteing these activities
	storePredecessorActivities, err := rt.store.GetPredecessorActivityNamesBySuccessorId(dto.ActivityId, r.Context())
	storePredecessorActivitiesStr := make([]string, len(storePredecessorActivities))
	for i, activity := range storePredecessorActivities {
		storePredecessorActivitiesStr[i] = activity.DispName
	}

	userActivitiesMap := make(map[string]struct{})
	for _, userActivity := range userInputs {
		userActivitiesMap[userActivity] = struct{}{}
	}
	userDeletedValues := valuesNotInMap(storePredecessorActivitiesStr, userActivitiesMap)
	for _, userDeletedValue := range userDeletedValues {
		predecessorActivity := storeActivitiesMap[userDeletedValue]
		if err := rt.store.DeleteDependencyByProjectPredecessorSuccessorId(dto.ProjectId, uuid.MustParse(predecessorActivity.ID), dto.ActivityId, r.Context()); err != nil {
			http.Error(w, "failed to delete activities from store", http.StatusInternalServerError)
			log.Printf("returned http internal server error while deleting activities from store. Encountered error: %v", err)
			return
		}
	}

	//display result
	rt.view.RenderDependencyInputResp(routemodels.HttpRenderTableRowDependencyInputResp{
		ProjectId:       dto.ProjectId,
		RowActivityId:   dto.ActivityId,
		ActivityNames:   userInputsNewValues,
		ResponseWriter:  w,
		Request:         r,
		DivTargetSwapId: fmt.Sprintf("dependency-suggestions-%s", dto.ActivityId),
		DivHxSwapOob:    "innerHTML",
	})
}

func removeDuplicates[T comparable](userInputs []T) []T {
	seen := make(map[T]struct{}, len(userInputs))
	n := 0

	for _, input := range userInputs {
		if _, exists := seen[input]; exists {
			continue
		}

		seen[input] = struct{}{}
		userInputs[n] = input
		n++
	}

	return userInputs[:n]
}

func removeEmptyString(userInputs []string) []string {
	seen := make(map[string]struct{}, len(userInputs))
	n := 0

	for _, input := range userInputs {
		if strings.TrimSpace(input) == "" {
			continue
		}

		seen[input] = struct{}{}
		userInputs[n] = input
		n++
	}

	return userInputs[:n]
}

func valuesNotInMap[K comparable, V any](values []K, valuesMap map[K]V) []K {
	unknownValues := []K{}
	for _, userInput := range values {
		if _, ok := valuesMap[userInput]; !ok {
			unknownValues = append(unknownValues, userInput)
		}
	}
	return unknownValues
}

func (rt *ViewRoutes) PostEmptyTableOfDependencyRowHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	var dto jsonmodels.CreateEmptyTableRow
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
		return
	}
	// calculate sort rank based on previous row and next row //
	defaultSortRankStep := int64(1000)
	sortRankStep := defaultSortRankStep
	sortRank := int64(0)

	if dto.PreviousRowSortRank != nil {
		sortRank = *dto.PreviousRowSortRank
	}

	if dto.NextRowSortRank != nil {
		sortRankStep = (*dto.NextRowSortRank - sortRank) / 2
	}

	sortRank = sortRank + sortRankStep
	////////////////////////////////////////////////////////////
	tableRow, err := rt.storeView.CreateEmptyDependencyTableRow(dto.ProjectId, sortRank, r.Context())
	if err != nil {
		http.Error(w, "failed to create new table row", http.StatusInternalServerError)
		log.Printf("returned http internal server error while creating empty dependency table")
		return
	}

	rt.view.RenderDependencyTableRow(routemodels.HttpRenderTableRow{
		TableRow:       tableRow,
		ResponseWriter: w,
		Request:        r,
	})
}

func (rt *ViewRoutes) PostTableOfDependencyRowDependencySuggestionSelectedHandler(w http.ResponseWriter, r *http.Request) {
	option := strings.TrimSpace(r.URL.Query().Get(string(routemodels.UrlQueryKeyOption)))
	if option == "" {
		http.Error(w, "option parameter missing from url", http.StatusBadRequest)
		return
	}
	jsonStr := r.FormValue(jsonmodels.JsonMarshalField)
	switch routemodels.DependencySuggestionType(option) {
	case routemodels.AddActivity:
		var dto jsonmodels.AddActivitiesFromRow
		err := json.Unmarshal([]byte(jsonStr), &dto)
		if err != nil {
			http.Error(w, "failed to parse json", http.StatusBadRequest)
			log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
			return
		}

		lastRow, err := rt.store.GetLastDependencyTableRowByProjectId(dto.IntoProjectId, r.Context())
		if err != nil {
			http.Error(w, "failed to get last row of table", http.StatusInternalServerError)
			log.Printf("returned http internal server error while getting last table row. encountered error: %v", err)
			return
		}
		newTableRows := []viewmodels.TableRow{}
		for i, activityName := range dto.ActivityNamesToAdd {
			sortRank := lastRow.SortRank + (viewmodels.TableRowSortRankStep * (int64(i) + 1))
			newTableRow, err := rt.storeView.CreateDependencyTableRow(dto.IntoProjectId, activityName, sortRank, r.Context())
			if err != nil {
				http.Error(w, "failed to create new table row", http.StatusInternalServerError)
				log.Printf("returned http internal server error while creating dependency table row")
				return
			}
			newTableRows = append(newTableRows, newTableRow)
		}
		rt.view.RenderDependencyInputSelectedAddActivityResp(routemodels.HttpRenderTableRowDependencyAddActivitySelectedResp{
			FromRowActivityId: dto.FromRowActivityId,
			TableRowsToAppend: newTableRows,
			ResponseWriter:    w,
			Request:           r,
		})
	default:
		http.Error(w, "unknown option specified", http.StatusBadRequest)
		return
	}
}
