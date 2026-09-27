package models

import "net/http"

type HttpParams struct {
	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request
}
