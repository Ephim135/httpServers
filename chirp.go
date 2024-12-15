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
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}
	return chirps
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	var authorID uuid.UUID
	var err error
	if author_id != "" {
		authorID, err = uuid.Parse(author_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "failed Parse id string to uuid", nil)
			return
		}
	}

	sort := r.URL.Query().Get("sort")

	// asc with author GetChirpsByAuthor
	if (sort == "asc" || sort == "") && author_id != "" {
		chirps, err := cfg.db.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Get Chirps sorted asc by created at Failed", nil)
			return
		}
		respondWithJSON(w, http.StatusOK, MapDatabaseChirps(chirps))
		return
	}
	// desc no author GetChirpsDesc
	if sort == "desc" && author_id != "" {
		chirps, err := cfg.db.GetChirpsByAuthorDesc(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Get Chirps sorted desc by created at Failed", nil)
			return
		}
		respondWithJSON(w, http.StatusOK, MapDatabaseChirps(chirps))
		return
	}
	// desc with author GetChirpsByAuthorDesc
	if sort == "desc" && author_id == "" {
		chirps, err := cfg.db.GetChirpsDesc(r.Context())
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Get Chirps sorted desc by created at Failed", nil)
			return
		}
		respondWithJSON(w, http.StatusOK, MapDatabaseChirps(chirps))
		return
	}

	// if no Paramaeters are passed or sort is asc
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Get Chirps Failed", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, MapDatabaseChirps(chirps))
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed Parse id string to uuid", nil)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid cant get Token:", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT invalid cant get Token:", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed Delete Chirp by ID", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
