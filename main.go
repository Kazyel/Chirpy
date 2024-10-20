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

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := api.CreateApiConfig(db, platform)
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Admin Routes
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	// API Routes
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerCreateUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerCreateChirps)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
