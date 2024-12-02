package auth_test

import (
	"testing"

	"github.com/Ephim135/httpServers.git/internal/auth"
)

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword"

	// Test that hashing the password does not return an error
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	// Test that the hashed password is not empty
	if hashedPassword == "" {
		t.Fatal("Hashed password is empty")
	}

	// Test that the hashed password is not the same as the original password
	if hashedPassword == password {
		t.Fatal("Hashed password should not match the plain text password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecretPassword"

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	// Test that CheckPasswordHash returns no error for a correct password
	err = auth.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error for the correct password: %v", err)
	}

	// Test that CheckPasswordHash returns an error for an incorrect password
	wrongPassword := "wrongPassword"
	err = auth.CheckPasswordHash(wrongPassword, hashedPassword)
	if err == nil {
		t.Fatal("CheckPasswordHash did not return an error for the wrong password")
	}
}
