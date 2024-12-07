package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	CreateNotification(ctx context.Context, notification *models.Notification) (*models.Notification, error)
	UpdateNotification(ctx context.Context, notification *models.Notification) (*models.Notification, error)
	GetNotifications(ctx context.Context, request dto.GetNotifications) ([]*models.Notification, error)
	GetNotification(ctx context.Context, id uint) (*models.Notification, error)
}

type NotificationRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewNotificationRepository(db db.DbService, logger logging.Logger) *NotificationRepository {
	return &NotificationRepository{db: db, logger: logger}
}
func (n *NotificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	err := n.db.GetDb().WithContext(ctx).Create(&notification).Error
	if err != nil {
		n.logger.Error(logging.Database, logging.Insert, "create notification error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return notification, err

}

func (n *NotificationRepository) GetNotification(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification

	err := n.db.GetDb().WithContext(ctx).First(&notification, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &notification, err
}

func (n *NotificationRepository) UpdateNotification(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	err := n.db.GetDb().WithContext(ctx).Save(notification).Error
	if err != nil {
		n.logger.Error(logging.Database, logging.Update, "Get choice by text and question id error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return notification, err
}
func (n *NotificationRepository) GetNotifications(ctx context.Context, request dto.GetNotifications) (notifications []*models.Notification, err error) {

	query := n.db.GetDb().WithContext(ctx)

	if request.UserId > 0 {
		query = query.Where("user_id = ?", request.UserId)
	}
	if request.Seen != nil {
		query = query.Where("seen = ?", *request.Seen)
	}

	if request.Sort != nil {
		query = query.Order(request.Sort.Field + " " + string(request.Sort.Type))
	}
	if request.Limit > 0 {
		query = query.Limit(int(request.Limit))
	}
	query = query.Offset(int(request.Offset))

	if err := query.Find(&notifications).Error; err != nil {
		n.logger.Error(logging.Database, logging.Select, "Get notifications error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

		return []*models.Notification{}, err
	}

	return notifications, nil
}
