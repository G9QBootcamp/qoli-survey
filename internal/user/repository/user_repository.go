package repository

import (
	"errors"
	"qoli-survey/internal/user/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(userID uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type userRepository struct {
	db DbService
}

func NewRepository(db DbService) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.GetDb().Create(user).Error
}

func (r *userRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.GetDb().First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.GetDb().Save(user).Error
}
