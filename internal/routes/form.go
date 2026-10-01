package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func (rt *FormRoutes) PutActivityDurationHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	var dto jsonmodels.UpdateActivityFromInput
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
		return
	}
	userInput := r.FormValue(dto.DomName)
	parseUserInput, err := parseDuration(userInput)
	if err != nil {
		http.Error(w, "failed to parse user input duration", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing user input: <%s>", userInput)
		return
	}
	if err := rt.store.UpdateActivityDuration(dto.ActivityId, parseUserInput, r.Context()); err != nil {
		http.Error(w, "failed to update activity name", http.StatusInternalServerError)
		log.Printf("returned http internal server error for err: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func parseDuration(s string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)d`)

	s = re.ReplaceAllStringFunc(s, func(part string) string {
		days, _ := strconv.ParseFloat(strings.TrimSuffix(part, "d"), 64)
		return fmt.Sprintf("%gh", days*24)
	})

	return time.ParseDuration(s)
}
