package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"uuid"

	"github.com/gtantech/p2/internal/models/jsonmodels"
	"github.com/gtantech/p2/internal/models/routemodels"
)

type HttpView interface {
	RenderHome(params routemodels.HttpHome)
	RenderDependencyTableRow(params routemodels.HttpRenderTableRow)
}

type ViewRoutes struct {
	view  HttpView
	store StoreView
}

func NewViewRoutes(view HttpView, store StoreView) *ViewRoutes {
	return &ViewRoutes{view: view, store: store}
}

func (rt *ViewRoutes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	mockProjectId := uuid.Max()
	table, err := rt.store.GetDependencyTableByProjectId(mockProjectId, r.Context())
	if err != nil {
		http.Error(w, "failed to get table", http.StatusInternalServerError)
	}
	rt.view.RenderHome(routemodels.HttpHome{Table: &table, ResponseWriter: w, Request: r})
}

func (rt *ViewRoutes) PostTableOfDependencyRowDependencyHandler(w http.ResponseWriter, r *http.Request) {

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
	tableRow, err := rt.store.CreateEmptyDependencyTableRow(dto.ProjectId, sortRank, r.Context())
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
