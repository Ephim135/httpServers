package main

import (
	"net/http"

	"github.com/Ephim135/httpServers.git/internal/auth"
)

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant get token from Bearer for revoking", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant revoke Token", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
