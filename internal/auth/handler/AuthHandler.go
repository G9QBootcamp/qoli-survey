package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/auth/service"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

func Signup(repo repository.UserRepository, jwtService service.JWTService) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(SignupRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Hash the password in the handler
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		// Create user
		user := models.User{
			Username: req.Username,
			Password: string(hashedPassword),
		}
		if err := repo.CreateUser(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		// Generate JWT token
		token, err := jwtService.GenerateToken(user.ID, user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, SignupResponse{Token: token})
	}
}
