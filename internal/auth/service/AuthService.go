package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	SecretKey string
	Issuer    string
}

type Claims struct {
	UserID     uint   `json:"user_id"`
	NationalID string `json:"national_id"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService initializes a new JWT service.
func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		SecretKey: secretKey,
		Issuer:    issuer,
	}
}

// GenerateToken generates a JWT token for a user.
func (s *JWTService) GenerateToken(userID uint, nationalID, email string) (string, error) {
	claims := &Claims{
		UserID:     userID,
		NationalID: nationalID,
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}
