package http

import (
	"fmt"
	"net/http"
	"strings"

	"log"

	"github.com/gorilla/mux"
	"github.com/pedrocmart/leaderboard-service/models"
)

func NewHandlersService(core *models.Core) models.HandlersService {
	return &BasicHandlers{core: core}
}

type BasicHandlers struct {
	core *models.Core
}

func ConnectBasic(router *mux.Router, core *models.Core) error {
	if router == nil {
		return fmt.Errorf("Could not connect default http mux handlers since router is nil")
	}
	basicAPI := BasicHandlers{core: core}
	router.HandleFunc("/user/{user_id}/score", basicAPI.HandleSubmitScore).Methods("POST")
	router.HandleFunc("/ranking", basicAPI.HandleGetRanking).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(basicAPI.NotFound)
	return nil
}

func (api *BasicHandlers) HandleSubmitScore(w http.ResponseWriter, r *http.Request) {
	if api.core.Service == nil {
		err := fmt.Errorf("Service is nil")
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	userId, ok := vars["user_id"]
	if !ok {
		err := fmt.Errorf("user_id is missing in parameters")
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}

	submitScoreRequest := new(models.SubmitScoreRequest)
	err := api.core.RequestResponse.ReadBodyAsJSON(r, submitScoreRequest)
	if err != nil {
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}

	result, err := api.core.Service.HandleSubmitScore(r.Context(), submitScoreRequest, userId)
	if err != nil {
		log.Printf("error while submiting score: %s", err.Error())
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}

	api.core.RequestResponse.HandleResponse(result, w, r, http.StatusOK)
}

func (api *BasicHandlers) HandleGetRanking(w http.ResponseWriter, r *http.Request) {
	if api.core.Service == nil {
		err := fmt.Errorf("Service is nil")
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}

	rankingType := r.URL.Query().Get("type")
	if strings.TrimSpace(rankingType) == "" {
		err := fmt.Errorf("A type must be included.")
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}

	result, err := api.core.Service.HandleGetRanking(r.Context(), rankingType)
	if err != nil {
		log.Printf("error while getting ranking: %s", err.Error())
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}

	api.core.RequestResponse.HandleResponse(result, w, r, http.StatusOK)
}

func (api *BasicHandlers) NotFound(w http.ResponseWriter, r *http.Request) {
	if api.core.Service == nil {
		err := fmt.Errorf("Service is nil")
		api.core.RequestResponse.HandleError(err, w, r, http.StatusInternalServerError)
		return
	}
	api.core.RequestResponse.HandleError(fmt.Errorf("404 not found"), w, r, http.StatusInternalServerError)
	return
}
