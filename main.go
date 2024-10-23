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
	// .env
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretToken := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	port := "8080"
	filepathRoot := "./app"
	api := api.CreateApiConfig(db, platform, secretToken, polkaKey)
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/app/", api.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Users Routes
	mux.HandleFunc("POST /api/users", api.HandlerCreateUsers)
	mux.Handle("PUT /api/users", api.MiddlewareAuthorize(http.HandlerFunc(api.HandlerUpdateUser)))

	// Chirps Routes
	mux.Handle("POST /api/chirps", api.MiddlewareAuthorize(http.HandlerFunc(api.HandlerCreateChirps)))
	mux.Handle("DELETE /api/chirps/{chirpID}", api.MiddlewareAuthorize(http.HandlerFunc(api.HandlerDeleteChirp)))
	mux.HandleFunc("GET /api/chirps", api.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", api.HandlerGetChirpByID)

	// Auth Routes
	mux.HandleFunc("POST /api/login", api.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", api.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", api.HandlerRevokeToken)

	// Admin Routes
	mux.HandleFunc("GET /admin/metrics", api.HandlerMetrics)
	mux.Handle("POST /admin/reset", api.MiddlewareDevMode(http.HandlerFunc(api.HandlerReset)))

	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", api.HandlerUpdateChirpy)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
