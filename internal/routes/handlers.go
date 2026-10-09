package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models"
	"github.com/gtantech/p2/internal/models/jsonmodels"
	"github.com/gtantech/p2/static"
)

type Routes struct {
	presenter models.Presenter
	store     models.Model
}

func NewRoutes(presenter models.Presenter, store models.Model) *Routes {
	r := &Routes{presenter: presenter, store: store}
	presenter.RegisterModel(store)
	return r
}

func (rt *Routes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeProjectId := uuid.Max() //mock home project id
	rt.presenter.DisplayHomeHandler(homeProjectId).ServeHTTP(w, r)
}

func (rt *Routes) GetHomeStyleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}

func (rt *Routes) PostFromRowPlusBtnReturnsEmptyTableRowHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	dto := jsonmodels.PostFromRowPlusBtn{}
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
		return
	}
	activity, err := rt.store.CreateActivity(dto.ProjectId, "", 0)
	if err != nil {
		http.Error(w, "failed to create new activity", http.StatusInternalServerError)
		log.Printf("returned http internal server error status while %v. Encountered error %v", "creating activity in store", err)
		return
	}
	trs, err := rt.store.GetTableRows(dto.ProjectId)
	if err != nil {
		http.Error(w, "failed to get table rows", http.StatusInternalServerError)
		log.Printf("returned http internal server error status while %v. Encountered error %v", "getting table rows from store", err)
		return
	}
	var rowAfter *models.TableRow = nil
	var rowCurrent models.TableRow
	for i := range trs {
		if trs[i].ActivityId == dto.RelativeToActivityId {
			rowCurrent = trs[i]
		}
		if i == 0 {
			continue
		}
		if trs[i-1].ActivityId == dto.RelativeToActivityId {
			rowAfter = &trs[i]
			break
		}
	}
	insertRowRank := rowCurrent.SortRank
	if rowAfter == nil {
		insertRowRank += 1000
	} else {
		insertRowRank = (rowCurrent.SortRank / 2) + (rowAfter.SortRank / 2)
	}

	newTr, err := rt.store.CreateTableRow(dto.ProjectId, activity.ActivityId, insertRowRank)
	if err != nil {
		http.Error(w, "failed to create new table row", http.StatusInternalServerError)
		log.Printf("returned http internal server error status while %v. Encountered error %v", "creating new table row in store", err)
		return
	}
	rt.presenter.DisplayTableRow(newTr).ServeHTTP(w, r)
}

func (rt *Routes) PutActivityNameHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	var dto jsonmodels.PutFromRowActivityNameChange
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: %s", jsonStr)
		return
	}
	userInput := r.FormValue(dto.DomName)
	if _, err := rt.store.UpdateActivityName(dto.ActivityId, userInput); err != nil {
		http.Error(w, "failed to update activity name", http.StatusInternalServerError)
		log.Printf("returned http internal server error for err: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *Routes) PutActivityDurationHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := r.FormValue("json")
	var dto jsonmodels.PutFromRowActivityDurationChange
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
		return
	}
	userInput := r.FormValue(dto.DomName)
	parseUserInput, err := parseDuration(userInput)
	if err != nil {
		http.Error(w, "failed to parse user input duration", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing user input: <%s>", userInput)
		return
	}
	if _, err := rt.store.UpdateActivityDuration(dto.ActivityId, parseUserInput); err != nil {
		http.Error(w, "failed to update activity name", http.StatusInternalServerError)
		log.Printf("returned http internal server error for err: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func parseDuration(s string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)d`)

	s = re.ReplaceAllStringFunc(s, func(part string) string {
		days, _ := strconv.ParseFloat(strings.TrimSuffix(part, "d"), 64)
		return fmt.Sprintf("%gh", days*24)
	})

	return time.ParseDuration(s)
}
