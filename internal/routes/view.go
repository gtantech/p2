package routes

import (
	"github.com/gtantech/p2/internal/models"
)

type HttpView interface {
	DisplayHome(params models.HttpViewHomeParams)
	DisplayTable(params models.HttpViewTableParams)
	DisplayEmptyTableRow(params models.HttpViewDisplayEmptyTableRowParams)
	DisplayDependencyWrapperWithNewActivitySuggestion(params models.HttpViewDisplayDependencyWrapperWithNewActivitySuggestion)
	DisplayDependencyWrapper(params models.HttpViewDisplayDependencyWrapper)
}
