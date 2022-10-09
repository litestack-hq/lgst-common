package jwt

import (
	"testing"
)

var secret = "secret"

func TestCreateToken(t *testing.T) {
	token, err := CreateToken("7038133d-3983-4684-b29d-78d46eb577d7", secret)

	if err != nil {
		t.Errorf("Failed to create JWT with error: %v", err)
	}

	if token == "" {
		t.Error("Created an empty JWT token")
	}
}

func TestVerifyToken(t *testing.T) {
	const userId = "7038133d-3983-4684-b29d-78d46eb577d7"
	validToken, _ := CreateToken(userId, secret)
	badToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cCI6MTUxNjIzOTAyMn0.ZPOgiA4zak2rtOLMEzATQ1N-Py238nxKZ-RBXhb7oe8"
	expToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1MTYyMzkwMjJ9.9RKzhmlwYbT9nJt8GAZxhhuvfIy4z1nWyTmSSCXxRxw"

	claim, err := VerifyToken(validToken, secret)

	if err != nil {
		t.Errorf("JWT verification failed with error: %v", err)
	}

	_, err = VerifyToken(validToken, "incorrect-secret")
	if err == nil {
		t.Error("JWT with the wrong key worked")
	}

	if claim == nil {
		t.Error("JWT verification failed")
	} else {
		if claim.UserID != userId {
			t.Error("JWT claim has incorrect data")
		}
	}

	_, err = VerifyToken(badToken, secret)

	if err == nil {
		t.Error("Invalid token type was unmarshalled")
	}

	_, err = VerifyToken(expToken, secret)
	if err == nil {
		t.Errorf("Verified expired token")
	}
}
