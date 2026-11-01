package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"uuid"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrRowNotFound     = errors.New("row not found")
)

type ProjectNotFoundError struct {
	ID uuid.UUID
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project not found with project id %v", e.ID)
}

func (e *ProjectNotFoundError) WriteToHttpResponse(w http.ResponseWriter) {
	http.Error(w, "project not found", http.StatusBadRequest)
}

type RowNotFoundError struct {
	ID uuid.UUID
}

func (e *RowNotFoundError) Error() string {
	return fmt.Sprintf("row not found with row id %v", e.ID)
}

func (e *RowNotFoundError) WriteToHttpResponse(w http.ResponseWriter) {
	http.Error(w, "row not found", http.StatusBadRequest)
}

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
