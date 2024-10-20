package api

import (
	"database/sql"
	"sync/atomic"

	"github.com/Kazyel/chirpy-bootdev/internal/database"
)

func CreateApiConfig(db *sql.DB, platform string) *ApiConfig {
	return &ApiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
	}
}
