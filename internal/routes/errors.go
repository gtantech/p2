package routes

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"uuid"

	"github.com/gtantech/p2/internal/logger"
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

// storeErrorHttpHandler handles store errors and writes to http response. err is assumed to be non-nil
func storeErrorHttpHandler(err error, w http.ResponseWriter) {
	if err, is := errors.AsType[*ProjectNotFoundError](err); is {
		err.WriteToHttpResponse(w)
		return
	}
	if err, is := errors.AsType[*RowNotFoundError](err); is {
		err.WriteToHttpResponse(w)
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

// errorToSlogAttr returns the slog.Attr with data from err.
// if the error type is unknown, internal server error with err value will be used for the return value.
func errorToSlogAttr(err error) slog.Attr {
	if err, is := errors.AsType[*ProjectNotFoundError](err); is {
		return logger.HttpErrorGroup(err.httpStatusCode(), err.httpErrMsg(), "projectId", err.ProjectId)
	}
	if err, is := errors.AsType[*RowNotFoundError](err); is {
		return logger.HttpErrorGroup(err.httpStatusCode(), err.httpErrMsg(), "rowId", err.RowId)
	}
	return logger.HttpErrorGroup(http.StatusInternalServerError, err.Error())
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
