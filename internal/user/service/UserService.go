package service

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
)

type UserService interface {
	Signup(c context.Context, password string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Signup(c context.Context, req dto.SignupRequest) error {

	user := models.User{
		City:         req.City,
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		NationalID:   req.NationalID,
		PasswordHash: req.PasswordHash,
	}

	// Check if email or national ID is already taken
	if s.repo.IsEmailOrNationalIDTaken(user.Email, user.NationalID) {
		return errors.New("email or national ID already in use")
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.PasswordHash = hashedPassword

	// Save the user
	return s.repo.CreateUser(user)
}
