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

type PutFromRowActivityDependencyChange struct {
	ActivityId uuid.UUID `json:"activityId"`
	ProjectId  uuid.UUID `json:"projectId"`
	DomName    string    `json:"domName"`
}

type PutFromRowActivityDurationChange struct {
	ActivityId uuid.UUID `json:"activityId"`
	DomName    string    `json:"domName"`
}

type PostFromDependencySuggestionAddNew struct {
	ActivityId                   uuid.UUID `json:"activityId"`
	ProjectId                    uuid.UUID `json:"projectId"`
	ActivitiesToCreate           []string  `json:"activitesToCreate"`
	ActivitiesToCreateFinalIndex []int     `json:"activitiesToCreateFinalIndex"`
}
