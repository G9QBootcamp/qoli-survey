package service

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"golang.org/x/net/context"
)

type IUserService interface {
	GetUsers(context.Context, dto.UserGetRequest) []*dto.UserResponse
	CreateUser(c context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error)
}
type UserService struct {
	conf *config.Config
	repo repository.IUserRepository
}

func New(conf *config.Config, repo repository.IUserRepository) *UserService {
	return &UserService{conf: conf, repo: repo}
}

func (s *UserService) GetUsers(c context.Context, r dto.UserGetRequest) []*dto.UserResponse {
	userFilters := dto.UserFilters{Name: r.Name}
	users, _ := s.repo.GetUsers(c, userFilters)
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
	return usersResponse
}

func (s *UserService) CreateUser(c context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error) {
	user := models.User{
		NationalID:   req.NationalID,
		Email:        req.Email,
		PasswordHash: util.HashPassword(req.Password),
	}

	user, err := s.repo.CreateUser(c, user)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:         user.ID,
		NationalID: user.NationalID,
		Email:      user.Email,
	}, nil
}
