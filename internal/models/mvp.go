package models

import (
	"net/http"
	"time"
	"uuid"
)

type Presenter interface {
	RegisterModel(model Model)
	DisplayHomeHandler(projectId uuid.UUID) http.Handler
}

type Model interface {
	GetActivities(projectId uuid.UUID) ([]Activity, error)
	CreateActivity(projectId uuid.UUID, activityName string, duration time.Duration) (Activity, error)
	GetDependencies(activityId uuid.UUID) ([]Activity, error)
}
