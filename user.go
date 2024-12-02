package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Ephim135/httpServers.git/internal/auth"
	"github.com/Ephim135/httpServers.git/internal/database"

	"github.com/google/uuid"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := userRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Fatalf("Failed to Hash Password %v", err)
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Fatalf("cant create user: %v", err)
		return
	}

	user := MapDatabaseUser(dbUser)

	respondWithJSON(w, http.StatusCreated, user)
}

func MapDatabaseUser(dbUser database.User) User {
	return User{
		ID:             dbUser.ID,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
	}
}
