package repository

import "github.com/G9QBootcamp/qoli-survey/internal/db"

type UserRepository struct {
	db db.DbService
}

func NewUserRepository(db db.DbService) *UserRepository {
	return &UserRepository{db: db}
}
