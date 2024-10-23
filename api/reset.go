package api

import "net/http"

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.db.DeleteUsers(r.Context())
}
