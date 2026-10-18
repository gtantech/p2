package routes

import (
	"errors"
	"log"
	"net/http"

	"github.com/gtantech/p2/internal/store"
)

// Returns true if the caller should return early
func storeErrorHandler(err error, w http.ResponseWriter, internalServerErrorLogMessage string) bool {
	if errors.Is(err, store.ErrProjectNotFound) {
		http.Error(w, "project not found", http.StatusBadRequest)
		return true
	}
	if errors.Is(err, store.ErrActivityNotFound) {
		http.Error(w, "activity not found", http.StatusBadRequest)
		return true
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
	log.Printf("returned http internal server error. %s", internalServerErrorLogMessage)
	return true
}
