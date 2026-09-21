package view

import (
	"time"
	"uuid"
)

type activity struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}
