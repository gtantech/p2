package jsonmodels

import "uuid"

type PostFromRowPlusBtn struct {
	ProjectId            uuid.UUID `json:"projectId"`
	RelativeToActivityId uuid.UUID `json:"relativeToActivityId"`
}
