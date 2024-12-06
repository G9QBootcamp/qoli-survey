package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/auth/models"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateOTP(context.Context, models.OTP) error
	GetOTPByUserID(context.Context, uint) (*models.OTP, error)
	UpdateOTP(context.Context, *models.OTP) error
}

type AuthRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewAuthRepository(db db.DbService, logger logging.Logger) *AuthRepository {
	return &AuthRepository{db: db, logger: logger}
}

func (r *AuthRepository) CreateOTP(c context.Context, otp models.OTP) error {
	err := r.db.GetDb().Create(&otp).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create otp error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *AuthRepository) GetOTPByUserID(c context.Context, userID uint) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.GetDb().
		Where("user_id = ? AND is_valid = ?", userID, true).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error(logging.Database, logging.Select, "Get otp error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *AuthRepository) UpdateOTP(c context.Context, otp *models.OTP) error {
	err := r.db.GetDb().Save(&otp).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "update otp error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}
