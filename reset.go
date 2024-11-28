package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, http.StatusForbidden, "only developers are allowed this endpoint", nil)
	}

	err := cfg.db.DelteAllUser(r.Context())
	if err != nil {
		log.Fatalf("Failed to delete users from table")
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}
