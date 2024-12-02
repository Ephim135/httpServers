package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ephim135/httpServers.git/internal/auth"
)

type loginRequest struct {
	Password string
	Email    string
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Fatalf("cant get User by Email: %v", err)
		return
	}

	user := MapDatabaseUser(dbUser)

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong Password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
