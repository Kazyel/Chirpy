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
}

func CreateApiConfig(db *sql.DB, platform string) *ApiConfig {
	return &ApiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
	}
}
