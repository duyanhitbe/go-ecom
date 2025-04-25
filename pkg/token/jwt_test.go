package token

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwt_Sign(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		claims     *Claims
		expectSign bool
	}{
		{
			name:       "Invalid secret",
			secret:     "very short",
			claims:     NewClaims("test", 0),
			expectSign: false,
		},
		{
			name:       "Nil claims",
			secret:     faker.Password() + faker.Username(),
			claims:     nil,
			expectSign: false,
		},
		{
			name:       "Generate token successfully",
			secret:     faker.Password() + faker.Username(),
			claims:     NewClaims("test", 0),
			expectSign: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jwt, err := NewJwt(tc.secret)
			fmt.Println(err)

			if tc.expectSign {
				tk, signErr := jwt.Sign(tc.claims)
				require.NoError(t, signErr)
				require.NotEmpty(t, tk)
			} else {
				if len(tc.secret) < 32 {
					require.Error(t, err)
				} else {
					_, signErr := jwt.Sign(tc.claims)
					require.Error(t, signErr)
				}
			}
		})
	}
}

func TestJwt_Verify(t *testing.T) {
	tests := []struct {
		name        string
		secret      string
		expectValid bool
	}{
		{
			name:        "Invalid token",
			secret:      faker.Password() + faker.Username(),
			expectValid: false,
		},
		{
			name:        "Valid token",
			secret:      faker.Password() + faker.Username(),
			expectValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jwt, _ := NewJwt(tc.secret)

			if tc.expectValid {
				sub := "test"
				claims := NewClaims(sub, time.Hour)
				tk, err := jwt.Sign(claims)
				require.NoError(t, err)

				verifiedClaims, err := jwt.Verify(tk)
				require.NotEmpty(t, verifiedClaims)
				require.Equal(t, sub, verifiedClaims.Subject)
				require.Equal(t, time.Hour, verifiedClaims.ExpiresAt.Sub(verifiedClaims.IssuedAt.Time))
			} else {
				verifiedClaims, err := jwt.Verify("invalid token")
				require.Error(t, err)
				require.Nil(t, verifiedClaims)
			}
		})
	}
}
