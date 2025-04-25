package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	secret []byte
}

func NewJwt(secret string) (*Jwt, error) {
	if len(secret) < 32 {
		return nil, errors.New("secret must be at least 32 characters long")
	}

	return &Jwt{
		secret: []byte(secret),
	}, nil
}

func (j *Jwt) Sign(payload *Claims) (string, error) {
	if payload == nil {
		return "", errors.New("payload cannot be nil")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(j.secret)
}

func (j *Jwt) Verify(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
