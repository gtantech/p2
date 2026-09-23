package view

import (
	"time"
	"uuid"
)

type Activity struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

func NewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *Activity {
	return &Activity{ID: id, ProjectID: projectId, DisplayName: displayName, Duration: duration}
}
