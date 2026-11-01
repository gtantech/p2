package routes

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"uuid"
)

type ProjectNotFoundError struct {
	ProjectId uuid.UUID
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project not found with project id %v", e.ProjectId)
}

func (e *ProjectNotFoundError) httpErrMsg() string {
	return "project not found"
}

func (e *ProjectNotFoundError) httpStatusCode() int {
	return http.StatusBadRequest
}

func (e *ProjectNotFoundError) WriteToHttpResponse(w http.ResponseWriter) {
	http.Error(w, e.httpErrMsg(), e.httpStatusCode())
}

func (e *ProjectNotFoundError) Slog() slog.Attr {
	return slog.Group("http", "status", e.httpStatusCode(), "err", e.httpErrMsg(), "projectId", e.ProjectId)
}

type RowNotFoundError struct {
	RowId uuid.UUID
}

func (e *RowNotFoundError) Error() string {
	return fmt.Sprintf("row not found with row id %v", e.RowId)
}

func (e *RowNotFoundError) httpErrMsg() string {
	return "row not found"
}

func (e *RowNotFoundError) httpStatusCode() int {
	return http.StatusBadRequest
}

func (e *RowNotFoundError) WriteToHttpResponse(w http.ResponseWriter) {
	http.Error(w, e.httpErrMsg(), e.httpStatusCode())
}

func (e *RowNotFoundError) Slog() slog.Attr {
	return slog.Group("http", "status", e.httpStatusCode(), "err", e.httpErrMsg(), "rowId", e.RowId)
}

// storeErrorHandler handles store errors and writes to http response. err is assumed to be non-nil. Returns true if the caller should return early
func storeErrorHandler(err error, w http.ResponseWriter, internalServerErrorLogMessage string) bool {
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
