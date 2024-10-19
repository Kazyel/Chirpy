package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func profaneFilter(s string) string {
	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	splittedString := strings.Split(s, " ")

	for _, word := range profaneWords {
		for i, v := range splittedString {
			if strings.ToLower(v) == word {
				splittedString[i] = "****"
			}
		}
	}

	joinedString := strings.Join(splittedString, " ")
	return joinedString
}

func respondWithJSON(w http.ResponseWriter, statusCode int, str string) {
	type jsonClean struct {
		CleanedBody string `json:"cleaned_body"`
	}

	filteredString := profaneFilter(str)
	cleanedJson := jsonClean{
		CleanedBody: filteredString,
	}
	cleanedMarshal, _ := json.Marshal(cleanedJson)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(cleanedMarshal)
}

func respondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
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

func handlerValidateChirpy(w http.ResponseWriter, r *http.Request) {
	type jsonRequest struct {
		Body string `json:"body"`
	}

	req := &jsonRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, 400, "Chirpy is too long")
	}

	respondWithJSON(w, 200, req.Body)
}
