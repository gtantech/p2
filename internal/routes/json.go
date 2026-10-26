package routes

import (
	"net/http"

	"github.com/gtantech/p2/internal/models/jsonmodels"
)

func getJson(r *http.Request) string {
	return r.FormValue(jsonmodels.JsonMarshalField)
}
