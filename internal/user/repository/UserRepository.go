package repository

import (
	"context"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
)

type IUserRepository interface {
	GetUsers(context.Context, dto.UserFilters) []*models.User
}

type UserRepository struct {
	db db.DbService
}

func NewUserRepository(db db.DbService) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(c context.Context, f dto.UserFilters) []*models.User {
	user := &models.User{Id: 1, Name: "aa"}
	var users []*models.User = []*models.User{}
	users = append(users, user)

	return users

}
