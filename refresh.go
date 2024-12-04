package main

import (
	"net/http"
	"time"

	"github.com/Ephim135/httpServers.git/internal/auth"
)

type tokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant get token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant get user for token from database", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant make new Token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tokenResponse{
		Token: token,
	},
	)
}
