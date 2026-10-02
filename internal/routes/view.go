package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"uuid"

	"github.com/gtantech/p2/internal/models/jsonmodels"
	"github.com/gtantech/p2/internal/models/routemodels"
	"github.com/gtantech/p2/internal/models/storemodels"
)

type HttpView interface {
	RenderHome(params routemodels.HttpHome)
	RenderDependencyTableRow(params routemodels.HttpRenderTableRow)
	RenderDependencySuggestion(params routemodels.HttpRenderTableRowDependencyActivitySuggestion)
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

	//reconcile user input with storeActivities:
	// - by returning a suggestion for creating unknown activities,
	// - deleting activities from depenendency store where no longer in user input
	storeActivities, err := rt.store.GetActivitiesByProjectId(dto.ProjectId, r.Context())
	storeActivitiesMap := make(map[string]storemodels.Activity)
	for _, storeActivity := range storeActivities {
		storeActivitiesMap[storeActivity.DispName] = storeActivity
	}

	//determine new values (unknown activities)
	userInputsNewValues := valuesNotInMap(userInputs, storeActivitiesMap)
	if len(userInputsNewValues) > 0 {
		rt.view.RenderDependencySuggestion(routemodels.HttpRenderTableRowDependencyActivitySuggestion{
			ActivityNames:  userInputsNewValues,
			RowActivityId:  dto.ActivityId,
			ResponseWriter: w,
			Request:        r,
		})
	}

	w.WriteHeader(http.StatusOK)
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
