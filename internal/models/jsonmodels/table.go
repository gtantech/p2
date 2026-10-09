package jsonmodels

import "uuid"

type PostFromRowPlusBtn struct {
	ProjectId            uuid.UUID `json:"projectId"`
	RelativeToActivityId uuid.UUID `json:"relativeToActivityId"`
}

type PutFromRowActivityNameChange struct {
	ActivityId uuid.UUID `json:"activityId"`
	DomName    string    `json:"domName"`
}
