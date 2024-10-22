package utils

import (
	"strings"
)

func ProfaneFilter(s string) string {
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
