package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/pedrocmart/leaderboard-service/coreservices"
	"github.com/pedrocmart/leaderboard-service/mocks"
	"github.com/pedrocmart/leaderboard-service/models"
	"github.com/stretchr/testify/assert"
)

func TestNewHandlersService(t *testing.T) {
	cases := []struct {
		description    string
		core           *models.Core
		expectedResult models.HandlersService
	}{
		{
			description: "should return a new handler",
			core:        &models.Core{},
			expectedResult: &BasicHandlers{
				core: &models.Core{},
			},
		},
	}
	for _, tc := range cases {
		result := NewHandlersService(tc.core)
		assert.Equal(t, tc.expectedResult, result, tc.description)
	}

}

func TestConnectBasic(t *testing.T) {
	cases := []struct {
		description    string
		req            string
		url            string
		expectedResult models.HandlersService
		routerNil      bool
	}{
		{
			description: "should return error if router nill",
			routerNil:   true,
		},
		{
			description: "should return users",
			req:         "POST",
			url:         "/users/{user_id}/score",
			expectedResult: &BasicHandlers{
				core: &models.Core{},
			},
		},
		{
			description: "should return ranking",
			req:         "GET",
			url:         "/ranking",
			expectedResult: &BasicHandlers{
				core: &models.Core{},
			},
		},
	}

	core := &models.Core{}
	core.Service = coreservices.NewCoreService(core)
	for _, tc := range cases {
		core := &models.Core{}
		router := mux.NewRouter()
		if tc.routerNil {
			router = nil
		}
		err := ConnectBasic(router, core)
		if tc.routerNil && err == nil {
			t.Errorf("Router null check not working")
		}
	}
}

func TestHandleSubmitScore(t *testing.T) {
	cases := []struct {
		description         string
		basicHandlers       BasicHandlers
		submitScoreResponse *models.SubmitScoreResponse
		submitScoreError    error
		queryParams         string
		core                *models.Core
		expectedError       string
		expectedResult      string
		requestBody         string
		service             bool
		writer              *httptest.ResponseRecorder
		request             *http.Request
		expectedStatusCode  int
		readJsonError       error
		vars                map[string]string
	}{
		{
			description:        "should creates user and submit score",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusOK,
			submitScoreResponse: &models.SubmitScoreResponse{
				UserID: 1,
				Score:  100,
			},
			requestBody: `{"score":"+100"}`,
			queryParams: `1`,
			service:     true,
			writer:      httptest.NewRecorder(),
			vars:        map[string]string{"user_id": "1"},
			request:     httptest.NewRequest("POST", "/user/user_id/score", bytes.NewReader([]byte(`{"score":"+100"}`))),
		},
		{
			description:        "should fail since the service is nil",
			expectedStatusCode: http.StatusInternalServerError,
			core:               &models.Core{},
			service:            false,
			writer:             httptest.NewRecorder(),
			request:            httptest.NewRequest("POST", "/user/user_id/score", bytes.NewReader([]byte(`{"score":"+100"}`))),
		},
		{
			description:        "should return an error whilst trying to parse body",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusInternalServerError,
			submitScoreResponse: &models.SubmitScoreResponse{
				UserID: 1,
				Score:  100,
			},
			readJsonError: fmt.Errorf("mock-parse-json-error"),
			requestBody:   `{"score":"+100"}`,
			queryParams:   `1`,
			service:       true,
			writer:        httptest.NewRecorder(),
			vars:          map[string]string{"user_id": "1"},
			request:       httptest.NewRequest("POST", "/user/user_id/score", bytes.NewReader([]byte(`{"score":"+100", "total":100}`))),
		},
		{
			description:        "should return an error when not sending an user_id",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusInternalServerError,
			submitScoreResponse: &models.SubmitScoreResponse{
				UserID: 1,
				Score:  100,
			},
			requestBody: `{"score":"+100"}`,
			queryParams: `1`,
			service:     true,
			writer:      httptest.NewRecorder(),
			request:     httptest.NewRequest("POST", "/user/user_id/score", bytes.NewReader([]byte(`{"score":"+100"}`))),
		},
		{
			description:        "should return an error when submiting score",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusInternalServerError,
			submitScoreError:   fmt.Errorf("mock-submitScore-error"),
			requestBody:        `{"score":"+100"}`,
			queryParams:        `1`,
			service:            true,
			vars:               map[string]string{"user_id": "1"},
			writer:             httptest.NewRecorder(),
			request:            httptest.NewRequest("POST", "/user/user_id/score", bytes.NewReader([]byte(`{"score":"+100"}`))),
		},
	}

	for _, tc := range cases {
		mockedService := mocks.ServiceMock{
			HandleSubmitScoreFunc: func(contextMoqParam context.Context, submitScoreRequest *models.SubmitScoreRequest, s string) (*models.SubmitScoreResponse, error) {
				return tc.submitScoreResponse, tc.submitScoreError
			},
		}
		if tc.service {
			tc.core.Service = &mockedService
		}

		requestResponseService := mocks.RequestResponseMock{
			HandleErrorFunc: func(err error, w http.ResponseWriter, r *http.Request, status int) {
				w.WriteHeader(status)
			},
			HandleResponseFunc: func(body interface{}, w http.ResponseWriter, r *http.Request, status int) {
				w.WriteHeader(status)
			},
			ReadBodyAsJSONFunc: func(req *http.Request, dest interface{}) error {
				return tc.readJsonError
			},
		}
		tc.request = mux.SetURLVars(tc.request, tc.vars)
		tc.core.RequestResponse = &requestResponseService
		tc.basicHandlers.core = tc.core
		tc.basicHandlers.HandleSubmitScore(tc.writer, tc.request)
		assert.Equal(t, tc.expectedStatusCode, tc.writer.Code, tc.description)
	}
}

func TestHandleGetRanking(t *testing.T) {
	cases := []struct {
		description        string
		basicHandlers      BasicHandlers
		getRankingResponse *models.GetRankingResponse
		getRankingError    error
		core               *models.Core
		expectedError      string
		expectedResult     string
		service            bool
		writer             *httptest.ResponseRecorder
		request            *http.Request
		expectedStatusCode int
		readJsonError      error
		vars               map[string]string
	}{
		{
			description:        "should gets ranking",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusOK,
			getRankingResponse: &models.GetRankingResponse{},
			service:            true,
			writer:             httptest.NewRecorder(),
			vars:               map[string]string{"user_id": "1"},
			request:            httptest.NewRequest("GET", "/score?type=top100", nil),
		},
		{
			description:        "should fail since the service is nil",
			expectedStatusCode: http.StatusInternalServerError,
			core:               &models.Core{},
			service:            false,
			writer:             httptest.NewRecorder(),
			request:            httptest.NewRequest("GET", "/score", nil),
		},
		{
			description:        "should return an error without a type",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusInternalServerError,
			getRankingResponse: &models.GetRankingResponse{},
			service:            true,
			writer:             httptest.NewRecorder(),
			vars:               map[string]string{"user_id": "1"},
			request:            httptest.NewRequest("GET", "/score", nil),
		},
		{
			description:        "should return an error getting the ranking",
			expectedResult:     `{"userid":1,"score":100}`,
			core:               &models.Core{},
			expectedStatusCode: http.StatusInternalServerError,
			getRankingError:    fmt.Errorf("mock-submitScore-error"),
			service:            true,
			vars:               map[string]string{"user_id": "1"},
			writer:             httptest.NewRecorder(),
			request:            httptest.NewRequest("GET", "/score?type=top100", nil),
		},
	}

	for _, tc := range cases {
		mockedService := mocks.ServiceMock{
			HandleGetRankingFunc: func(contextMoqParam context.Context, s string) (*models.GetRankingResponse, error) {
				return tc.getRankingResponse, tc.getRankingError
			},
		}
		if tc.service {
			tc.core.Service = &mockedService
		}

		requestResponseService := mocks.RequestResponseMock{
			HandleErrorFunc: func(err error, w http.ResponseWriter, r *http.Request, status int) {
				w.WriteHeader(status)
			},
			HandleResponseFunc: func(body interface{}, w http.ResponseWriter, r *http.Request, status int) {
				w.WriteHeader(status)
			},
		}
		tc.request = mux.SetURLVars(tc.request, tc.vars)
		tc.core.RequestResponse = &requestResponseService
		tc.basicHandlers.core = tc.core
		tc.basicHandlers.HandleGetRanking(tc.writer, tc.request)
		assert.Equal(t, tc.expectedStatusCode, tc.writer.Code, tc.description)
	}
}

func TestNotFound(t *testing.T) {
	cases := []struct {
		description                  string
		basicHandlers                BasicHandlers
		core                         *models.Core
		writer                       *httptest.ResponseRecorder
		request                      *http.Request
		params                       httprouter.Params
		expectedStatusCode           int
		expectedNotFoundError        error
		service                      bool
		expectedRequestResponseError error
		requestResponseWriter        http.ResponseWriter
	}{
		{
			description:        "should fail since the service is nil",
			params:             nil,
			expectedStatusCode: http.StatusInternalServerError,
			core:               &models.Core{},
			service:            false,
			writer:             httptest.NewRecorder(),
			request:            httptest.NewRequest("GET", "/healthz", nil),
		},
		{
			description:           "should return error if version function errors out",
			params:                nil,
			expectedStatusCode:    http.StatusInternalServerError,
			core:                  &models.Core{},
			service:               true,
			writer:                httptest.NewRecorder(),
			request:               httptest.NewRequest("GET", "/healthz", nil),
			expectedNotFoundError: errors.New("error"),
		},
	}
	for _, tc := range cases {
		mockedService := mocks.ServiceMock{}
		if tc.service {
			tc.core.Service = &mockedService
		}

		requestResponseService := mocks.RequestResponseMock{
			HandleErrorFunc: func(err error, w http.ResponseWriter, r *http.Request, status int) {
				w.WriteHeader(status)
			},
			HandleResponseFunc: func(body interface{}, w http.ResponseWriter, r *http.Request, status int) {
				w.WriteHeader(status)
			},
		}
		tc.core.RequestResponse = &requestResponseService
		tc.basicHandlers.core = tc.core
		tc.basicHandlers.NotFound(tc.writer, tc.request)
		assert.Equal(t, tc.expectedStatusCode, tc.writer.Code, tc.description)
	}
}
