package service

import (
	"errors"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"github.com/G9QBootcamp/qoli-survey/pkg/jwtutils"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IUserService interface {
	GetUsers(context.Context, dto.UserGetRequest) ([]*dto.UserResponse, error)
	Login(c context.Context, req dto.LoginRequest) (string, time.Time, error)
	UpdateUserProfile(c context.Context, userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateUserNotifications(userID uint, req *dto.UpdateNotificationsRequest) (*models.User, error)
}
type UserService struct {
	conf   *config.Config
	repo   repository.IUserRepository
	logger logging.Logger
}

func New(conf *config.Config, repo repository.IUserRepository, logger logging.Logger) *UserService {
	return &UserService{conf: conf, repo: repo, logger: logger}
}

func (s *UserService) GetUsers(c context.Context, r dto.UserGetRequest) ([]*dto.UserResponse, error) {
	userFilters := dto.UserFilters{Name: r.Name}
	users, err := s.repo.GetUsers(c, userFilters)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGetUsers, "error in get users", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return nil, err
	}
	usersResponse := []*dto.UserResponse{}

	for _, user := range users {
		usersResponse = append(usersResponse, ToUserResponse(&user))
	}
	return usersResponse, nil
}

func (s *UserService) UpdateUserProfile(c context.Context, userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(c, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth) // Assuming format "YYYY-MM-DD"
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		if time.Since(user.CreatedAt) > 24*time.Hour {
			return nil, errors.New("date of birth cannot be updated after 24 hours of registration")
		}
		user.DateOfBirth = dateOfBirth
	}
	if req.City != "" {
		user.City = req.City
	}

	updatedUser, err := s.repo.UpdateUser(c, user)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToUpdateUser, "error in updating user", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}

	return ToUserResponse(updatedUser), nil
}

func ToUserResponse(user *models.User) *dto.UserResponse {
	var dateOfBirth string

	if !user.DateOfBirth.IsZero() {
		formattedDate := user.DateOfBirth.Format("2006-01-02")
		dateOfBirth = formattedDate
	}

	return &dto.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: dateOfBirth,
		City:        user.City,
	}
}

func (s *UserService) UpdateUserNotifications(userID uint, req *dto.UpdateNotificationsRequest) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.SurveyCanceled != nil {
		user.Notifications.SurveyCanceled = *req.SurveyCanceled
	}
	if req.VoteCanceled != nil {
		user.Notifications.VoteCanceled = *req.VoteCanceled
	}
	if req.RoleAssigned != nil {
		user.Notifications.RoleAssigned = *req.RoleAssigned
	}

	return s.repo.Update(user)
}
func (s *UserService) Login(c context.Context, req dto.LoginRequest) (string, time.Time, error) {
	user, err := s.repo.GetUserByEmail(c, req.Email)
	if err != nil {
		s.logger.Error(logging.Internal, logging.UserNotAuthorized, "invalid email", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return "", time.Time{}, errors.New("invalid email")
	}

	if !user.EmailVerified {
		s.logger.Info(logging.Internal, logging.UserNotAuthorized, "email is not verified", map[logging.ExtraKey]interface{}{logging.Service: "UserService"})
		return "", time.Time{}, errors.New("email is not verified")
	}

	if err := util.CheckPassword(req.Password, user.PasswordHash); err != nil {
		s.logger.Error(logging.Internal, logging.UserNotAuthorized, "invalid email", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return "", time.Time{}, errors.New("invalid password")
	}

	expiresAt := time.Now().Add(time.Duration(s.conf.JWT.ExpireMinutes) * time.Minute)
	token, err := jwtutils.GenerateToken(user.ID, user.GlobalRole.Name, s.conf.JWT.SecretKey, s.conf.JWT.ExpireMinutes)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGenerateToken, "generate token error", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return "", time.Time{}, errors.New("failed to generate token")
	}

	return token, expiresAt, nil
}
