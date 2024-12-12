package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid request", err)
		return
	}

	if req.Event != "user.created" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid user id", err)
		return
	}

	if req.Event == "user.created" {
		cfg.db.UpdateUserRed(r.Context(), userID)
		w.WriteHeader(http.StatusNoContent)
	}
}
