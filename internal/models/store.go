package models

import (
	"time"
	"uuid"
)

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

type StoreDependency struct {
	ID                    uuid.UUID
	ProjectID             uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

type StoreGetPredecessorNamesBySuccessorResult struct {
	DependencyID            uuid.UUID
	Relationship            RelationshipType
	PredecessorActivityID   uuid.UUID
	PredecessorActivityName string
}

type StoreCreateDepdencencyParams struct {
	ProjectID             uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

type StoreUpdateDepdencencyParams struct {
	ID                    uuid.UUID
	Relationship          RelationshipType
	PredecessorActivityID uuid.UUID
	SuccessorActivityID   uuid.UUID
}

type StoreActivity struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type StoreCreateActivityParams struct {
	ProjectID   uuid.UUID
	DisplayName string
	Duration    time.Duration
}

type StoreGetActivityByNameAndProjectParams struct {
	ProjectID   uuid.UUID
	DisplayName string
}

type StoreUpdateActivityParams struct {
	Id          uuid.UUID
	DisplayName string
	Duration    time.Duration
}
