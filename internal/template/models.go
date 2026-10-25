package template

import (
	"strings"

	"github.com/gtantech/p2/internal/models"
)

type tableRow struct {
	ActivityName string
	Dependencies string
	Duration     string
}

func ToTableRow(params models.TableRow) *tableRow {
	return &tableRow{
		ActivityName: params.ActivityName,
		Dependencies: func() string {
			return strings.Join(params.DependencyNames, ", ")
		}(),
		Duration: params.Duration.String(),
	}
}
