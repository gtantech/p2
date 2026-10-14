package models

import "uuid"

type Project struct {
	ProjectId          uuid.UUID
	ProjectDisplayName string
}
