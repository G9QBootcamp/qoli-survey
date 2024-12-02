package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUsers(ctx context.Context, filters dto.UserFilters) ([]models.User, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type UserRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewUserRepository(db db.DbService, logger logging.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
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
	query = query.Where("deleted_at is null")

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	if err := query.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []models.User{}, nil
		}
		r.logger.Error(logging.Database, logging.Select, "get users error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

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
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create user error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return user, err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	return r.db.GetDb().WithContext(ctx).Where("ID = ?", id).Delete(&models.User{}).Error
}
