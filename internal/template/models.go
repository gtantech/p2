package template

import (
	"strings"
	"uuid"

	"github.com/gtantech/p2/internal/models"
)

type tableRow struct {
	RowId        uuid.UUID
	ActivityName string
	Dependencies string
	Duration     string
}

func ToTableRow(params models.TableRow) *tableRow {
	return &tableRow{
		RowId:        params.RowId,
		ActivityName: params.ActivityName,
		Dependencies: func() string {
			return strings.Join(params.DependencyNames, ", ")
		}(),
		Duration: params.Duration.String(),
	}
}
