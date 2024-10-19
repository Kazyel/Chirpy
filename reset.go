package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("This endpoint is only available in dev mode."))
	}

	cfg.db.DeleteUsers(r.Context())
}
