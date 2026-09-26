package models

import "uuid"

type StoreProject struct {
	ID          uuid.UUID
	DisplayName string
}

type StoreCreateProjectParams struct {
	DisplayName string
}

type StoreUpdateProjectParams struct {
	Id          uuid.UUID
	DisplayName string
}
