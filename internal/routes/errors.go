package routes

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrRowNotFound     = errors.New("row not found")
)

// Returns true if the caller should return early
func storeErrorHandler(err error, w http.ResponseWriter, internalServerErrorLogMessage string) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrProjectNotFound) {
		http.Error(w, "project not found", http.StatusBadRequest)
		return true
	}
	if errors.Is(err, ErrRowNotFound) {
		http.Error(w, "row not found", http.StatusBadRequest)
		return true
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
	log.Printf("returned http internal server error. %s", internalServerErrorLogMessage)
	return true
}

// Returns true if the caller should return early
func jsonUnmarshalErrorHandler(err error, w http.ResponseWriter, jsonStr string) bool {
	if err == nil {
		return false
	}
	http.Error(w, "failed to parse json", http.StatusBadRequest)
	log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
	return true
}
