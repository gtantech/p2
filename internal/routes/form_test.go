package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/gtantech/p2/internal/models/jsonmodels"
)

type mockStore struct {
	UpdateActivityDurationCallback func(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error
	UpdateActivityNameCallback     func(activityId uuid.UUID, activityName string, ctx context.Context) error
}

// UpdateActivityDuration implements [Store].
func (m *mockStore) UpdateActivityDuration(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error {
	return m.UpdateActivityDurationCallback(activityId, activityDuration, ctx)
}

// UpdateActivityName implements [Store].
func (m *mockStore) UpdateActivityName(activityId uuid.UUID, activityName string, ctx context.Context) error {
	return m.UpdateActivityNameCallback(activityId, activityName, ctx)
}

var _ Store = (*mockStore)(nil) //ensures mockStore implements Store at compile time
func TestParseDuration(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"4d", 4 * 24 * time.Hour},
		{"2d3h", 2*24*time.Hour + 3*time.Hour},
		{"1d2h30m", 24*time.Hour + 2*time.Hour + 30*time.Minute},
	}

	for _, tt := range tests {
		got, err := parseDuration(tt.input)
		if err != nil {
			t.Fatalf("parseDuration(%q) returned error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("parseDuration(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestPutActivityDurationHandlerNoError(t *testing.T) {
	//set up mock client-sent data
	domName := "test-dom-name"
	dto := jsonmodels.UpdateActivityFromInput{ActivityId: uuid.New(), DomName: domName}
	formVals := url.Values{}
	formVals.Add("json", jsonmodels.MarshalParams(dto))
	mockUserInput := "2h"
	expectedDuration := 2 * time.Hour
	formVals.Add(domName, mockUserInput)

	//setup mock store to be called by handler
	store := mockStore{}
	store.UpdateActivityDurationCallback = func(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error {
		if got, want := activityId, dto.ActivityId; got != want {
			return fmt.Errorf("got %v, want %v", got, want)
		}
		if got, want := activityDuration, expectedDuration; got != want {
			return fmt.Errorf("got %v, want %v", got, want)
		}

		return nil
	}

	//register handler
	route := NewFormRoutes(&store)
	server := httptest.NewServer(http.HandlerFunc(route.PutActivityDurationHandler))

	//create request
	req, err := http.NewRequest(
		http.MethodPut,
		server.URL+"/",
		strings.NewReader(formVals.Encode()),
	)
	if err != nil {
		t.Fatal(err)
	}

	//set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	//check empty body
	if string(body) != `` {
		t.Errorf("unexpected response: %s", body)
	}

}

func TestPutActivityDurationHandlerStoreError(t *testing.T) {
	//set up mock client-sent data
	domName := "test-dom-name"
	dto := jsonmodels.UpdateActivityFromInput{ActivityId: uuid.New(), DomName: domName}
	formVals := url.Values{}
	formVals.Add("json", jsonmodels.MarshalParams(dto))
	mockUserInput := "2h"
	formVals.Add(domName, mockUserInput)

	//setup mock store to be called by handler
	store := mockStore{}
	mockStoreErrorMsg := "mocked error"
	store.UpdateActivityDurationCallback = func(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error {
		return fmt.Errorf("%s", mockStoreErrorMsg)
	}

	//register handler
	route := NewFormRoutes(&store)
	server := httptest.NewServer(http.HandlerFunc(route.PutActivityDurationHandler))

	//create request
	req, err := http.NewRequest(
		http.MethodPut,
		server.URL+"/",
		strings.NewReader(formVals.Encode()),
	)
	if err != nil {
		t.Fatal(err)
	}

	//set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := strings.TrimSpace(string(body)), "failed to update activity name"; got != want {
		t.Errorf("unexpected response: %s, want %s", got, want)
	}

}

func TestPutActivityDurationHandlerBadInput(t *testing.T) {
	//set up mock client-sent data
	domName := "test-dom-name"
	dto := jsonmodels.UpdateActivityFromInput{ActivityId: uuid.New(), DomName: domName}
	formVals := url.Values{}
	formVals.Add("json", jsonmodels.MarshalParams(dto))
	mockUserBadInput := "2h-3"
	formVals.Add(domName, mockUserBadInput)

	//setup mock store to be called by handler
	store := mockStore{}
	store.UpdateActivityDurationCallback = func(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error {
		t.Fatal("Unexpected call to store")
		return nil
	}

	//register handler
	route := NewFormRoutes(&store)
	server := httptest.NewServer(http.HandlerFunc(route.PutActivityDurationHandler))

	//create request
	req, err := http.NewRequest(
		http.MethodPut,
		server.URL+"/",
		strings.NewReader(formVals.Encode()),
	)
	if err != nil {
		t.Fatal(err)
	}

	//set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := strings.TrimSpace(string(body)), "failed to parse user input duration"; got != want {
		t.Errorf("unexpected response: %s, want %s", got, want)
	}

}

func TestPutActivityDurationHandlerBadJson(t *testing.T) {
	//set up mock client-sent data
	domName := "test-dom-name"
	formVals := url.Values{}
	formVals.Add("json", "bad-json-input")
	mockUserInput := "2h"
	formVals.Add(domName, mockUserInput)

	//setup mock store to be called by handler
	store := mockStore{}
	store.UpdateActivityDurationCallback = func(activityId uuid.UUID, activityDuration time.Duration, ctx context.Context) error {
		t.Fatal("Unexpected call to store")
		return nil
	}

	//register handler
	route := NewFormRoutes(&store)
	server := httptest.NewServer(http.HandlerFunc(route.PutActivityDurationHandler))

	//create request
	req, err := http.NewRequest(
		http.MethodPut,
		server.URL+"/",
		strings.NewReader(formVals.Encode()),
	)
	if err != nil {
		t.Fatal(err)
	}

	//set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := strings.TrimSpace(string(body)), "failed to parse json"; got != want {
		t.Errorf("unexpected response: %s, want %s", got, want)
	}

}

func TestPutActivityNameHandler(t *testing.T) {
	mockActivityId := uuid.New()
	mockActivityDispName := "test-activity-name"
	domName := "test-dom-name"
	tests := []struct {
		testName                   string
		formValueKey               string
		jsonStr                    string
		domKey                     string
		userInput                  string
		expectedHttpStatusCode     int
		expectedBody               string
		UpdateActivityNameCallback func(activityId uuid.UUID, activityName string, ctx context.Context) error
	}{
		{
			testName:               "NoError",
			formValueKey:           "json",
			jsonStr:                jsonmodels.MarshalParams(jsonmodels.UpdateActivityFromInput{ActivityId: mockActivityId, DomName: domName}),
			domKey:                 domName,
			userInput:              mockActivityDispName,
			expectedHttpStatusCode: http.StatusOK,
			expectedBody:           "",
			UpdateActivityNameCallback: func(activityId uuid.UUID, activityName string, ctx context.Context) error {
				if got, want := activityId, mockActivityId; got != want {
					return fmt.Errorf("got %v, want %v", got, want)
				}
				if got, want := activityName, mockActivityDispName; got != want {
					return fmt.Errorf("got %v, want %v", got, want)
				}

				return nil
			},
		},
		{
			testName:               "BadJson",
			formValueKey:           "json",
			jsonStr:                "badjsonstr",
			domKey:                 domName,
			userInput:              mockActivityDispName,
			expectedHttpStatusCode: http.StatusBadRequest,
			expectedBody:           "failed to parse json",
			UpdateActivityNameCallback: func(activityId uuid.UUID, activityName string, ctx context.Context) error {
				t.Fatal("unexpected store call")

				return nil
			},
		},
		{
			testName:               "StoreError",
			formValueKey:           "json",
			jsonStr:                jsonmodels.MarshalParams(jsonmodels.UpdateActivityFromInput{ActivityId: mockActivityId, DomName: domName}),
			domKey:                 domName,
			userInput:              mockActivityDispName,
			expectedHttpStatusCode: http.StatusInternalServerError,
			expectedBody:           "failed to update activity name",
			UpdateActivityNameCallback: func(activityId uuid.UUID, activityName string, ctx context.Context) error {
				if got, want := activityName, mockActivityDispName; got != want {
					return fmt.Errorf("got %v, want %v", got, want)
				}
				if got, want := activityId, uuid.Max(); got != want {
					return fmt.Errorf("got %v, want %v", got, want)
				}

				return nil
			},
		},
		{
			testName:               "BadJsonKey",
			formValueKey:           "badjsonkey",
			jsonStr:                jsonmodels.MarshalParams(jsonmodels.UpdateActivityFromInput{ActivityId: mockActivityId, DomName: domName}),
			domKey:                 domName,
			userInput:              mockActivityDispName,
			expectedHttpStatusCode: http.StatusBadRequest,
			expectedBody:           "failed to parse json",
			UpdateActivityNameCallback: func(activityId uuid.UUID, activityName string, ctx context.Context) error {
				t.Fatal("unexpected store call")

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//setup mock store to be called by handler
			store := mockStore{}
			store.UpdateActivityNameCallback = tt.UpdateActivityNameCallback

			//register handler
			route := NewFormRoutes(&store)
			server := httptest.NewServer(http.HandlerFunc(route.PutActivityNameHandler))

			//mock form sent from hx-put
			formVals := url.Values{}
			formVals.Add(tt.formValueKey, tt.jsonStr)
			formVals.Add(tt.domKey, tt.userInput)

			//create request
			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/",
				strings.NewReader(formVals.Encode()),
			)
			if err != nil {
				t.Fatal(err)
			}

			//set header
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			//make request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedHttpStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedHttpStatusCode, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			//check body
			if strings.TrimSpace(string(body)) != tt.expectedBody {
				t.Errorf("unexpected response: %s", body)
			}
		})
	}
}
