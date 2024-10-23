package api

import (
	"database/sql"
	"sync/atomic"

	"github.com/Kazyel/chirpy-bootdev/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secretToken    string
	polkaKey       string
}

func CreateApiConfig(db *sql.DB, platform string, token string, polkaKey string) *ApiConfig {
	return &ApiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
		secretToken:    token,
		polkaKey:       polkaKey,
	}
}
