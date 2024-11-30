package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUsers(ctx context.Context, filters dto.UserFilters) ([]models.User, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type UserRepository struct {
	db db.DbService
}

func NewUserRepository(db db.DbService) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(ctx context.Context, filters dto.UserFilters) ([]models.User, error) {
	var users []models.User

	query := r.db.GetDb().WithContext(ctx).Model(&models.User{})

	if filters.Name != "" {
		query = query.Where("CONCAT(first_name, ' ', last_name) ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filters.Email+"%")
	}

	if filters.NationalID != "" {
		query = query.Where("national_id ILIKE ?", "%"+filters.NationalID+"%")
	}

	if filters.City != "" {
		query = query.Where("city ILIKE ?", "%"+filters.City+"%")
	}

	if filters.YearOfBirth > 0 {
		query = query.Where("EXTRACT(YEAR FROM date_of_birth) = ?", filters.YearOfBirth)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User

	err := r.db.GetDb().WithContext(ctx).First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := r.db.GetDb().WithContext(ctx).Create(&user).Error
	return user, err
}
