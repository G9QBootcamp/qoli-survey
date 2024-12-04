package service

import (
	"errors"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IUserService interface {
	GetUsers(context.Context, dto.UserGetRequest) ([]*dto.UserResponse, error)
	Signup(c context.Context, req dto.SignupRequest) (*dto.UserResponse, error)
	UpdateUserProfile(c context.Context, userID uint, req *dto.UpdateUserRequest) (*models.User, error) 
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
		usersResponse = append(usersResponse, &dto.UserResponse{
			ID:          user.ID,
			NationalID:  user.NationalID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			City:        user.City,
			DateOfBirth: user.DateOfBirth,
		})
	}
	return usersResponse, nil
}

func (s *UserService) Signup(c context.Context, req dto.SignupRequest) (*dto.UserResponse, error) {

	user := models.User{
		NationalID:   req.NationalID,
		Email:        req.Email,
		PasswordHash: req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		City:         req.City,
		DateOfBirth:  req.DateOfBirth,
	}

	if s.repo.IsEmailOrNationalIDTaken(c, user.Email, user.NationalID) {
		return nil, errors.New("email or national ID already in use")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.logger.Error(logging.Internal, logging.HashPassword, "failed to hash password", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}
	user.PasswordHash = hashedPassword

	user, err = s.repo.CreateUser(c, user)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToCreateUser, "error in create user", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}

	return &dto.UserResponse{
		ID:          user.ID,
		NationalID:  user.NationalID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		City:        user.City,
		DateOfBirth: user.DateOfBirth,
	}, nil
}

func (s *UserService) UpdateUserProfile(c context.Context, userID uint, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetUserByID(c, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.DateOfBirth != nil {
		dateOfBirth, err := time.Parse("2006-01-02", *req.DateOfBirth) //assuming the date format is "YYYY-MM-DD"
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		user.DateOfBirth = dateOfBirth
	}
	if req.City != nil {
		user.City = *req.City
	}

	if req.DateOfBirth != nil {
		if time.Since(user.CreatedAt) > 24*time.Hour {
			return nil, errors.New("date of birth cannot be updated after 24 hours of registration")
		}
	}

	return s.repo.UpdateUser(c, user)
}
