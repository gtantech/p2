package routes

import (
	"github.com/gtantech/p2/internal/models"
)

type HttpView interface {
	Home(params models.ViewHomeParams)
	DisplayEmptyTableRow(params models.ViewDisplayEmptyTableRowParams)
}
