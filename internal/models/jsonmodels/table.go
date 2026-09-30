package jsonmodels

import (
	"uuid"
)

type CreateEmptyTableRow struct {
	ProjectId uuid.UUID `json:"projectId"`
}
