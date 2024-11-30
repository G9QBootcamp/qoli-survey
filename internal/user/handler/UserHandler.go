package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IUserService
}

func NewHandler(conf *config.Config, db db.DbService) *UserHandler {
	return &UserHandler{conf: conf, db: db, service: service.New(conf, repository.NewUserRepository(db))}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users := h.service.GetUsers(c.Request().Context(), dto.UserGetRequest{Name: "aa", Page: 1})
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.UserCreateRequest

	// Bind the request body to the DTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Call the service layer
	user, err := h.service.CreateUser(c.Request().Context(), req)
	if err != nil {
		// Handle other errors as internal server errors
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return the created user
	return c.JSON(http.StatusCreated, user)
}
