package routemodels

import (
	"net/http"

	"github.com/gtantech/p2/internal/models/viewmodels"
)

type HttpParams struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HttpHome struct {
	viewmodels.Home
	HttpParams
}

type HttpRenderTableRow struct {
	viewmodels.TableRow
	NextRowSortRank *int64
	HttpParams
}
