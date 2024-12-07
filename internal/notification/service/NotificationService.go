package service

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/models"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"

	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type INotificationService interface {
	Notify(c context.Context, userId uint, message string) (*dto.NotificationResponse, error)
	Seen(c context.Context, id uint) error
	GetNotifications(c context.Context, userId uint, request dto.GetNotificationsRequest) ([]*dto.NotificationResponse, error)
}
type NotificationService struct {
	conf   *config.Config
	repo   repository.INotificationRepository
	logger logging.Logger
}

func New(conf *config.Config, repo repository.INotificationRepository, logger logging.Logger) *NotificationService {
	return &NotificationService{conf: conf, repo: repo, logger: logger}
}

func (n *NotificationService) Notify(c context.Context, userId uint, message string) (response *dto.NotificationResponse, err error) {
	notif, err := n.repo.CreateNotification(c, &models.Notification{UserID: userId, Message: message})

	if err != nil {
		n.logger.Error(logging.Internal, logging.FailedToSendNotify, "error in send notification", map[logging.ExtraKey]interface{}{logging.Service: "Notification", logging.ErrorMessage: err.Error()})
		return nil, err
	}

	return response, util.ConvertTypes(n.logger, notif, &response)

}
func (n *NotificationService) Seen(c context.Context, id uint) error {
	notif, err := n.repo.GetNotification(c, id)
	if err != nil {
		n.logger.Error(logging.Internal, logging.FailedToSeenNotifications, "error in get notification", map[logging.ExtraKey]interface{}{logging.Service: "Notification", logging.ErrorMessage: err.Error()})
		return err
	}

	notif.Seen = true

	_, err = n.repo.UpdateNotification(c, notif)
	if err != nil {
		n.logger.Error(logging.Internal, logging.FailedToSeenNotifications, "error in get notification", map[logging.ExtraKey]interface{}{logging.Service: "Notification", logging.ErrorMessage: err.Error()})
		return err
	}

	return nil
}
func (n *NotificationService) GetNotifications(c context.Context, userId uint, request dto.GetNotificationsRequest) (response []*dto.NotificationResponse, err error) {
	limit := 10
	offset := 0
	if request.Page > 0 {
		offset = limit * (request.Page - 1)
	}

	notifs, err := n.repo.GetNotifications(c, dto.GetNotifications{UserId: userId, Seen: &request.Seen, Limit: limit, Offset: offset, Sort: &dto.Sort{Type: dto.DescSort, Field: "created_at"}})
	if err != nil {
		n.logger.Error(logging.Internal, logging.FailedToSeenNotifications, "error in get notification", map[logging.ExtraKey]interface{}{logging.Service: "Notification", logging.ErrorMessage: err.Error()})
		return []*dto.NotificationResponse{}, err
	}

	return response, util.ConvertTypes(n.logger, notifs, &response)
}
