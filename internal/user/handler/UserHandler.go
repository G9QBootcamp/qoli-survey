package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IUserService
	logger  logging.Logger
}

func NewHandler(conf *config.Config, db db.DbService, logger logging.Logger) *UserHandler {
	return &UserHandler{conf: conf, db: db, service: service.New(conf, repository.NewUserRepository(db, logger), logger), logger: logger}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers(c.Request().Context(), dto.UserGetRequest{Name: "aa", Page: 1})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Signup(c echo.Context) error {
	var req dto.SignupRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	user, err := h.service.Signup(c.Request().Context(), req)
	if err != nil {
		// Handle other errors as internal server errors
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	userID = 1
	if !ok {
		//return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := c.Validate(&req); err != nil {
		var validationErrors map[string]string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	updatedUser, err := h.service.UpdateUserProfile(c.Request().Context(), userID, req)
	if err != nil {
		if err.Error() == "date of birth cannot be updated after 24 hours of registration" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to update profile"})
	}

	return c.JSON(http.StatusOK, updatedUser)
}
