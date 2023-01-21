package models

import (
	"encoding/json"
	"net/http"

	"log"
)

type BasicRequestResponse struct {
}

func (brr BasicRequestResponse) log(endpoint string, message interface{}) {
	log.Printf("%s - %+v", endpoint, message)
}
func (brr BasicRequestResponse) HandleError(err error, w http.ResponseWriter, r *http.Request, status int) {
	if err == nil {
		return //there is no error to write back
	}
	brr.log(r.URL.String(), err)
	writeError := writeJsonError(err, w, status)
	if writeError != nil {
		brr.log("Error writing bytes to the Response Writer:\n%s", err.Error())
	}
}

func (brr BasicRequestResponse) HandleResponse(body interface{}, w http.ResponseWriter, r *http.Request, status int) {
	brr.log(r.URL.String(), body)
	err := writeJson(body, w, status)
	if err != nil {
		brr.log("Error writing bytes to the Response Writer:\n%s", err.Error())
	}
}

func (brr BasicRequestResponse) ReadBodyAsJSON(req *http.Request, dest interface{}) (err error) {
	err = json.NewDecoder(req.Body).Decode(dest)
	return
}

func writeJson(body interface{}, w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(body)
	if err != nil {
		writeErrorErr := writeJsonError(err, w, http.StatusInternalServerError)
		if err != nil {
			return writeErrorErr
		}
		return err
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	return err
}

func writeJsonError(err error, w http.ResponseWriter, status int) error {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	body := Response{
		Message: errorMessage,
		Success: false,
		Data:    nil,
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	return err
}

type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
