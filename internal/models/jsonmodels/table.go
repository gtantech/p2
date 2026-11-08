package jsonmodels

import "uuid"

type PostFromRowPlusBtn struct {
	RelativeToRowId uuid.UUID `json:"relativeToRowId"`
}

type PutFromRowActivityNameChange struct {
	RowId   uuid.UUID `json:"rowId"`
	DomName string    `json:"domName"`
}
