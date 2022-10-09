package jwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CustomClaim - Custom payload type
type CustomClaim struct {
	UserID string `json:"user_id,omitempty"`
	jwt.StandardClaims
}

// CreateToken - Create JWT token
func CreateToken(userID string, secret string) (string, error) {
	now := time.Now()

	claims := CustomClaim{
		userID,
		jwt.StandardClaims{
			ExpiresAt: now.Add(15 * time.Minute).Unix(),
			IssuedAt:  now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// VerifyToken - Verifies JWT
func VerifyToken(tokenString string, key string) (*CustomClaim, error) {
	claims := CustomClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return &claims, nil
	}

	return nil, errors.New("failed to verify JWT")
}
