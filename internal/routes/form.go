package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gtantech/p2/internal/models/jsonmodels"
)

type FormRoutes struct {
	store Store
}

func NewFormRoutes(store Store) *FormRoutes {
	return &FormRoutes{store: store}
}

func (rt *FormRoutes) PutActivityNameHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	var dto jsonmodels.UpdateActivityFromInput
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: %s", jsonStr)
		return
	}
	userInput := r.FormValue(dto.DomName)
	if err := rt.store.UpdateActivityName(dto.ActivityId, userInput, r.Context()); err != nil {
		http.Error(w, "failed to update activity name", http.StatusInternalServerError)
		log.Printf("returned http internal server error for err: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
