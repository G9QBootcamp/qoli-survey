package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	conf *config.Config
	db   db.DbService
	repo repository.UserRepository
}

func NewHandler(conf *config.Config, db db.DbService) *UserHandler {
	return &UserHandler{conf: conf, db: db, repo: *repository.NewUserRepository(db)}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Echo!")
}
