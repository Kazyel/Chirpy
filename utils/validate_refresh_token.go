package utils

import (
	"net/http"
	"time"
)

func ValidateRefreshToken(revoked bool, expires time.Time, w http.ResponseWriter) {
	if revoked {
		RespondWithError(w, 401, "Refresh token revoked")
		return
	}

	if expires.Before(time.Now()) {
		RespondWithError(w, 401, "Refresh token expired")
		return
	}
}
