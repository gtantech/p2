package routes

import (
	"github.com/gtantech/p2/internal/models"
)

type HttpView interface {
	DisplayHome(params models.HttpViewHomeParams)
	DisplayEmptyTableRow(params models.HttpViewDisplayEmptyTableRowParams)
}
