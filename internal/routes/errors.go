package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"uuid"
)

type ProjectNotFoundError struct {
	ProjectId uuid.UUID
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project not found with project id %v", e.ProjectId)
}

func (e *ProjectNotFoundError) WriteToHttpResponse(w http.ResponseWriter) {
	http.Error(w, "project not found", http.StatusBadRequest)
}

type RowNotFoundError struct {
	RowId uuid.UUID
}

func (e *RowNotFoundError) Error() string {
	return fmt.Sprintf("row not found with row id %v", e.RowId)
}

func (e *RowNotFoundError) WriteToHttpResponse(w http.ResponseWriter) {
	http.Error(w, "row not found", http.StatusBadRequest)
}

// Returns true if the caller should return early
func storeErrorHandler(err error, w http.ResponseWriter, internalServerErrorLogMessage string) bool {
	if err == nil {
		return false
	}
	if err, is := errors.AsType[*ProjectNotFoundError](err); is {
		err.WriteToHttpResponse(w)
		return true
	}
	if err, is := errors.AsType[*RowNotFoundError](err); is {
		err.WriteToHttpResponse(w)
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
