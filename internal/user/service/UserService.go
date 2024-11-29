package service

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"golang.org/x/net/context"
)

type IUserService interface {
	GetUsers(context.Context, dto.UserRequest) []*dto.UserResponse
}
type UserService struct {
	conf *config.Config
	repo repository.IUserRepository
}

func New(conf *config.Config, repo repository.IUserRepository) *UserService {
	return &UserService{conf: conf, repo: repo}
}

func (s *UserService) GetUsers(c context.Context, r dto.UserRequest) []*dto.UserResponse {
	userFilters := dto.UserFilters{Name: r.Name}
	users := s.repo.GetUsers(c, userFilters)
	usersResponse := []*dto.UserResponse{}

	for _, v := range users {
		usersResponse = append(usersResponse, &dto.UserResponse{Name: v.Name})
	}
	return usersResponse
}
