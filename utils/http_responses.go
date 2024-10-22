package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	type jsonError struct {
		Error string `json:"error"`
	}

	jsonErr := jsonError{
		Error: errorMessage,
	}

	errorResponse, err := json.Marshal(jsonErr)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(errorResponse)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, errorCode int, response interface{}) {
	marshalResponse, err := json.Marshal(response)

	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	if statusCode == 200 {
		w.Write(marshalResponse)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		return
	}

	w.WriteHeader(statusCode)
	w.Write(marshalResponse)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}
