package routemodels

import (
	"net/http"
	"uuid"

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

type HttpRenderTableRowDependencyInputResp struct {
	ProjectId     uuid.UUID
	RowActivityId uuid.UUID
	ActivityNames []string
	HttpParams
	DivTargetSwapId string
	DivHxSwapOob    string
}

type HttpRenderTableRowDependencyAddActivitySelectedResp struct {
	FromRowActivityId uuid.UUID
	TableRow          *viewmodels.TableRow
	HttpParams
}
