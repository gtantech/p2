package jsonmodels

import "uuid"

type PostFromRowPlusBtn struct {
	RelativeToRowId uuid.UUID `json:"relativeToRowId"`
}
