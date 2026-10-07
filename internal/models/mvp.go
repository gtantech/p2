package models

import (
	"net/http"
	"uuid"
)

type Presenter interface {
	RegisterModel(model Model)
	DisplayHomeHandler(projectId uuid.UUID) http.Handler
}

type Model interface {
	GetActivities(projectId uuid.UUID) ([]Activity, error)
	GetDependencies(activityId uuid.UUID) ([]Activity, error)
}
