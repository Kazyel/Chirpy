package api

import (
	"net/http"

	"github.com/Kazyel/chirpy-bootdev/internal/auth"
	"github.com/Kazyel/chirpy-bootdev/utils"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareAuthorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := auth.GetBearerToken(r.Header)

		if err != nil {
			utils.RespondWithError(w, 401, err.Error())
			return
		}

		userID, err := auth.ValidateJWT(bearerToken, cfg.secretToken)

		if err != nil {
			utils.RespondWithError(w, 401, err.Error())
			return
		}

		r.Header.Set("User-ID", userID.String())
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareDevMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.platform != "dev" {
			w.WriteHeader(403)
			w.Write([]byte("This endpoint is only available in dev mode."))
			return
		}

		next.ServeHTTP(w, r)
	})
}
