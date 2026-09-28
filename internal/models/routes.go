package models

import (
	"time"
	"uuid"
)

type RoutesPostEmptyTableRow struct {
	ProjectId uuid.UUID `json:"projectId"`
}

type RoutesPostNewActivity struct {
	ProjectId   uuid.UUID     `json:"projectId"`
	DisplayName string        `json:"displayName"`
	Duration    time.Duration `json:"duration"`
}
