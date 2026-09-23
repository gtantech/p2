package view

import (
	"time"
	"uuid"
)

type Activity struct {
	id          uuid.UUID
	projectID   uuid.UUID
	displayName string
	duration    time.Duration
}

func NewActivity(id uuid.UUID, projectId uuid.UUID, displayName string, duration time.Duration) *Activity {
	return &Activity{id: id, projectID: projectId, displayName: displayName, duration: duration}
}
