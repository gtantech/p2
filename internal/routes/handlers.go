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
	return r
}

func (rt *Routes) HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeProjectId := uuid.Max() //mock home project id
	tableRows, err := rt.store.GetTableRowsSortedByRank(r.Context(), homeProjectId)
	if err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error status while %v. Encountered error %v", "getting sorted table rows from store", err)) {
			return
		}
	}
	table := models.Table{
		ProjectId: homeProjectId,
		Rows:      tableRows,
	}
	rt.presenter.HtmlHomeHandler(table).ServeHTTP(w, r)
}

func (rt *Routes) GetHomeStyleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(static.StaticHomeCss)
}

func (rt *Routes) PostFromRowPlusBtnReturnsEmptyTableRowHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := getJson(r)
	dto := jsonmodels.PostFromRowPlusBtn{}
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: <%s>", jsonStr)
		return
	}
	activity, err := rt.store.CreateActivity(r.Context(), dto.ProjectId, "", 0)
	if err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error status while %v. Encountered error %v", "creating activity in store", err)) {
			return
		}
	}
	newTr, err := rt.store.CreateTableRowBelow(r.Context(), dto.ProjectId, activity.ActivityId, dto.RelativeToActivityId)
	if err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error status while %v. Encountered error %v", "creating new table row in store", err)) {
			return
		}
	}
	rt.presenter.HtmlTableRow(newTr).ServeHTTP(w, r)
}

func (rt *Routes) PutActivityNameHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := getJson(r)
	var dto jsonmodels.PutFromRowActivityNameChange
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: %s", jsonStr)
		return
	}
	userInput := r.FormValue(dto.DomName)
	if _, err := rt.store.UpdateActivityName(r.Context(), dto.ActivityId, userInput); err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error while updating activity name from store. Encountered error: %s", err)) {
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *Routes) PostActivityDependencySuggestionAddNewHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := getJson(r)
	var dto jsonmodels.PostFromDependencySuggestionAddNew
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: %s", jsonStr)
		return
	}
	createdTableRows := []models.TableRow{}
	for _, name := range dto.ActivitiesToCreate {
		a, err := rt.store.CreateActivity(r.Context(), dto.ProjectId, name, 0)
		if err != nil {
			if storeErrorHandler(err, w, fmt.Sprintf("error while creating new activity: %s", name)) {
				return
			}
		}
		tr, err := rt.store.CreateTableRowAtBottom(r.Context(), dto.ProjectId, a.ActivityId)
		if err != nil {
			if storeErrorHandler(err, w, fmt.Sprintf("error while creating new table row: %s", name)) {
				return
			}
		}
		createdTableRows = append(createdTableRows, tr)
	}
	dependencies, err := rt.store.GetDependenciesBySuccessorActivityId(r.Context(), dto.ActivityId)
	if err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error while getting dependencies by successor id with successor id: %v. Encountered error: %v", dto.ActivityId, err)) {
			return
		}
	}
	createdDependencyIdxMap := make(map[int]uuid.UUID)
	for i := range len(dto.ActivitiesToCreateFinalIndex) {
		idx := dto.ActivitiesToCreateFinalIndex[i]
		createdDependencyIdxMap[idx] = createdTableRows[i].ActivityId
	}
	dependencyCount := len(dependencies) + len(dto.ActivitiesToCreate)
	knownDependencyIdx := 0
	updatedDependencyIds := make([]uuid.UUID, dependencyCount)
	for i := range dependencyCount {
		if _, ok := createdDependencyIdxMap[i]; !ok {
			updatedDependencyIds[i] = dependencies[knownDependencyIdx].PredecessorActivityId
			knownDependencyIdx += 1
		} else {
			updatedDependencyIds[i] = createdDependencyIdxMap[i]
		}
	}

	if _, err := rt.store.UpdateActivityDependencies(r.Context(), dto.ProjectId, dto.ActivityId, updatedDependencyIds); err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error while updating dependency. Encountered error %v", err)) {
			return
		}
	}

	rt.presenter.HtmxPostActivityDependencySuggestionAddNewHandlerResp(createdTableRows, dto.ActivityId).ServeHTTP(w, r)
}

func (rt *Routes) PostActivityDependencyHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := getJson(r)
	var dto jsonmodels.PutFromRowActivityDependencyChange
	err := json.Unmarshal([]byte(jsonStr), &dto)
	if err != nil {
		http.Error(w, "failed to parse json", http.StatusBadRequest)
		log.Printf("returned http bad request error while parsing json: %s", jsonStr)
		return
	}

	activities, err := rt.store.GetActivities(r.Context(), dto.ProjectId)
	if err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error while getting activities from store. Encountered error: %s", err)) {
			return
		}
	}

	activityNameMap := make(map[string]uuid.UUID)

	for _, activity := range activities {
		activityNameMap[activity.DisplayName] = activity.ActivityId
	}
	dependencyIds := []uuid.UUID{}

	userInputs := strings.Split(r.FormValue(dto.DomName), ",")
	userAddNewActiivtyName := []string{}
	unknownDependencyFinalIdx := []int{}
	for i := range userInputs {
		userInput := strings.TrimSpace(userInputs[i])
		if userInputId, ok := activityNameMap[userInput]; ok {
			dependencyIds = append(dependencyIds, userInputId)
		} else if userInput != "" {
			userAddNewActiivtyName = append(userAddNewActiivtyName, userInput)
			unknownDependencyFinalIdx = append(unknownDependencyFinalIdx, i)
		}
	}

	if _, err := rt.store.UpdateActivityDependencies(r.Context(), dto.ProjectId, dto.ActivityId, dependencyIds); err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error while updating activity dependencies from store. Encountered error: %s", err)) {
			return
		}
	}

	rt.presenter.HtmxPostActivityDependencyHandlerResp(userAddNewActiivtyName, unknownDependencyFinalIdx, dto.ActivityId, dto.ProjectId).ServeHTTP(w, r)
}

func (rt *Routes) PutActivityDurationHandler(w http.ResponseWriter, r *http.Request) {
	jsonStr := getJson(r)
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
	if _, err := rt.store.UpdateActivityDuration(r.Context(), dto.ActivityId, parseUserInput); err != nil {
		if storeErrorHandler(err, w, fmt.Sprintf("error while updating activity duration from store. Encountered error: %s", err)) {
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func getJson(r *http.Request) string {
	return r.FormValue(jsonmodels.JsonMarshalField)
}

func parseDuration(s string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)d`)

	s = re.ReplaceAllStringFunc(s, func(part string) string {
		days, _ := strconv.ParseFloat(strings.TrimSuffix(part, "d"), 64)
		return fmt.Sprintf("%gh", days*24)
	})

	return time.ParseDuration(s)
}
