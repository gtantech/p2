package models

import (
	"net/http"
)

type HttpParams struct {
	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request
}

type HttpViewDisplayEmptyTableRowParams struct {
	ViewDisplayEmptyTableRowParams
	HttpParams
}

type HttpViewHomeParams struct {
	ViewHomeParams
	HttpParams
}
