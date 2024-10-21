package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Kazyel/chirpy-bootdev/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	port := "8080"
	filepathRoot := "./app"

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretToken := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	api := api.CreateApiConfig(db, platform, secretToken)
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/app/", api.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Admin Routes
	mux.HandleFunc("GET /admin/metrics", api.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", api.HandlerReset)

	// API Routes
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", api.HandlerCreateUsers)

	mux.HandleFunc("POST /api/chirps", api.HandlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", api.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", api.HandlerGetChirpByID)
	mux.HandleFunc("POST /api/chirps/{userID}", api.HandlerGetChirpByUserID)

	// Auth Routes
	mux.HandleFunc("POST /api/login", api.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", api.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", api.HandlerRevokeToken)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
