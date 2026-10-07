package models

import (
	"time"
	"uuid"
)

type Activity struct {
	ActivityId  uuid.UUID
	DisplayName string
	Duration    time.Duration
}
