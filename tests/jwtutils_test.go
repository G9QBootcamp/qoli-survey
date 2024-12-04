package test

import (
	"testing"

	"github.com/G9QBootcamp/qoli-survey/pkg/jwtutils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	secretKey := "secret"
	userID := uint(1)
	role := "user"
	expireMinutes := 60

	token, err := jwtutils.GenerateToken(userID, role, secretKey, expireMinutes)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(userID), claims["user_id"])
	assert.Equal(t, role, claims["role"])
}
