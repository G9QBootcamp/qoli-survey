package service

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
)

type IUserService interface {
	GetUsers(dto.UserRequest) []*dto.UserResponse
}
type UserService struct {
	conf *config.Config
	repo repository.IUserRepository
}

func New(conf *config.Config, repo repository.IUserRepository) *UserService {
	return &UserService{conf: conf, repo: repo}
}

func (s *UserService) GetUsers(dto.UserRequest) []*dto.UserResponse {
	users := s.repo.GetUsers()
	usersResponse := []*dto.UserResponse{}

	for _, v := range users {
		usersResponse = append(usersResponse, &dto.UserResponse{Name: v.Name})
	}
	return usersResponse
}
