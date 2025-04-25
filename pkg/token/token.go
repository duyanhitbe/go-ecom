package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
}

func NewClaims(sub string, duration time.Duration) *Claims {
	now := time.Now()

	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   sub,
			Issuer:    sub,
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(duration)},
		},
	}
}

type Token interface {

	// Sign generates a signed JWT string from the provided Claims object and returns it along with any potential error.
	Sign(payload *Claims) (string, error)

	// Verify parses and validates a given token string, returning the token's claims if it is valid or an error if it is not.
	Verify(token string) (*Claims, error)
}
