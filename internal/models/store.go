package models

import "uuid"

type StoreProject struct {
	ID          uuid.UUID
	DisplayName string
}

type CreateProjectParams struct {
	DisplayName string
}

type UpdateProjectParams struct {
	Id          uuid.UUID
	DisplayName string
}
