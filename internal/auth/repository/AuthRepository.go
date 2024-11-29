package repository

import (
	"github.com/G9QBootcamp/qoli-survey/internal/auth/models"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
)

type IAuthRepository interface {
	StoreOTP() error
	GetValidOTP() (*models.OTP, error)
	InvalidateOTP() error
}

type AuthRepository struct {
	db db.DbService
}

func NewAuthRepository(db db.DbService) *AuthRepository {
	return &AuthRepository{db: db}
}
