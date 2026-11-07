package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gtantech/p2/internal/logger"
	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/models/jsonmodels"
	"github.com/gtantech/p2/internal/template"
	"github.com/gtantech/p2/static"
)

type Routes struct {
	store Store
	log   *slog.Logger
}

func New(store Store, log *slog.Logger) *Routes {
	r := &Routes{
		store: store,
		log:   log.WithGroup(packagePath),
	}
	return r
}

func (rt *Routes) PostFromRowPlusBtnReturnsEmptyTableRow(w http.ResponseWriter, r *http.Request) {
	jsonStr := getJson(r)
	dto := jsonmodels.PostFromRowPlusBtn{}
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		jsonUnmarshalHttpErrorHandler(w)
		rt.log.Error("failed to unmarshal json", logger.JsonGroup(jsonStr))
		return
	}
	tr, err := rt.store.CreateRowBelow(r.Context(), dto.RelativeToRowId)
	if err != nil {
		storeErrorHttpHandler(err, w)
		rt.log.Error(
			fmt.Sprintf("failed to create new row below row id %v", dto.RelativeToRowId),
			errorToSlogAttr(err),
		)
		return
	}
	rt.ViewTableRow(w, r, tr)
}

func (rt *Routes) ViewTableRow(w http.ResponseWriter, r *http.Request, props models.TableRow) {
	template.DependencyTableRow(template.ToTableRow(props)).Render(r.Context(), w)
}

func (rt *Routes) GetHome(w http.ResponseWriter, r *http.Request) {
	homeProject, err := rt.store.GetHomeProject(r.Context())
	if err != nil {
		storeErrorHttpHandler(err, w)
		rt.log.Error(
			"failed to get home project",
			errorToSlogAttr(err),
		)
		return
	}
	props, err := rt.store.GetTable(r.Context(), homeProject.ProjecId)
	if err != nil {
		storeErrorHttpHandler(err, w)
		rt.log.Error(
			"failed to get home table",
			errorToSlogAttr(err),
		)
		return
	}
	rt.ViewHome(w, r, props)
}

func (rt *Routes) ViewHome(w http.ResponseWriter, r *http.Request, props models.Table) {
	template.Home(props, "Project Planner").Render(r.Context(), w)
}

func (rt *Routes) GetHomeStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}
