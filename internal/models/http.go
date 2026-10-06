package models

import "net/http"

type HttpParams struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}
