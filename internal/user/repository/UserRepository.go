package repository

import (
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
)

type IUserRepository interface {
	GetUsers() []*models.User
}

type UserRepository struct {
	db db.DbService
}

func NewUserRepository(db db.DbService) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers() []*models.User {
	user := &models.User{Id: 1, Name: "aa"}
	var users []*models.User = []*models.User{}
	users = append(users, user)

	return users

}
