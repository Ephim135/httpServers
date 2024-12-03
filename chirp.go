package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ephim135/httpServers.git/internal/auth"
	"github.com/Ephim135/httpServers.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid cant get Token:", err)
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := ChirpsValidate(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userID,
		Body:   cleaned,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Create Chirp Failed", nil)
		return
	}

	chirp := MapDatabaseChirp(dbChirp)
	respondWithJSON(w, http.StatusCreated, chirp)
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func MapDatabaseChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

func MapDatabaseChirps(dbChirps []database.Chirp) []Chirp {
	mappedChirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		mappedChirps[i] = MapDatabaseChirp(dbChirp)
	}
	return mappedChirps
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Get Chirps Failed", nil)
		return
	}
	chirpsMapped := MapDatabaseChirps(chirps)
	respondWithJSON(w, http.StatusOK, chirpsMapped)
}

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed Parse id string to uuid", nil)
		return
	}

	databaseChirp, err := cfg.db.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed Get Chirp by ID", nil)
		return
	}
	chirp := MapDatabaseChirp(databaseChirp)

	respondWithJSON(w, http.StatusOK, chirp)

}
