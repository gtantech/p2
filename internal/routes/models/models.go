package models

import "uuid"

type PostEmptyTableRow struct {
	ProjectId uuid.UUID `json:"projectId"`
}
